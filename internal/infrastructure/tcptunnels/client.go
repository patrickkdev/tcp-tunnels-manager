package tcptunnels

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/patrickkdev/tcptunnelsmanager/internal/domain"
)

type Event struct {
	TunnelID int
	Level    domain.LogLevel
	Message  string
}

type Tunnel struct {
	domain.TunnelRow

	cmd *exec.Cmd

	ctx    context.Context
	cancel context.CancelFunc

	wg sync.WaitGroup

	Events chan Event

	MaxRetryBackoff time.Duration
}

func New(row domain.TunnelRow) *Tunnel {
	ctx, cancel := context.WithCancel(context.Background())

	return &Tunnel{
		TunnelRow: row,
		ctx:       ctx,
		cancel:    cancel,
		Events:    make(chan Event, 32),

		MaxRetryBackoff: 30 * time.Second,
	}
}

func (t *Tunnel) Start() error {
	t.wg.Add(1)
	go t.run()

	return nil
}

func (t *Tunnel) run() {
	defer t.wg.Done()

	backoff := time.Second

	for {
		select {
		case <-t.ctx.Done():
			return
		default:
		}

		cmd := exec.Command(
			"socat",
			"-d",
			"-d",
			fmt.Sprintf("TCP-LISTEN:%d,fork,reuseaddr", t.ListenPort),
			fmt.Sprintf("TCP:%s:%d", t.TargetHost, t.TargetPort),
		)

		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()

		if err := cmd.Start(); err != nil {
			t.emit(domain.LogLevelError, fmt.Sprintf("failed to start socat: %v", err))
			return
		}

		t.cmd = cmd

		var wg sync.WaitGroup
		wg.Add(2)

		go t.pipeLogs(stdout, &wg)
		go t.pipeLogs(stderr, &wg)

		err := cmd.Wait()

		wg.Wait()

		select {
		case <-t.ctx.Done():
			return
		case <-time.After(backoff):
		}

		// Exponential backoff up to 30 seconds
		if backoff < t.MaxRetryBackoff {
			backoff *= 2
		} else {
			backoff = t.MaxRetryBackoff
		}

		if err != nil {
			t.emit(domain.LogLevelError,
				fmt.Sprintf("socat crashed: %v, restarting", err))
		} else {
			t.emit(domain.LogLevelError,
				"socat exited unexpectedly, restarting")
		}
	}
}

func (t *Tunnel) pipeLogs(pipe io.ReadCloser, wg *sync.WaitGroup) {
	defer wg.Done()

	scanner := bufio.NewScanner(pipe)

	for scanner.Scan() {
		t.emit(domain.LogLevelInfo, scanner.Text())
	}
}

func (t *Tunnel) emit(level domain.LogLevel, msg string) {
	select {
	case t.Events <- Event{
		TunnelID: t.ID,
		Level:    level,
		Message:  msg,
	}:
	default:
	}
}

func (t *Tunnel) Stop() error {
	t.cancel()

	if t.cmd != nil && t.cmd.Process != nil {
		_ = t.cmd.Process.Kill()
	}

	t.wg.Wait()

	close(t.Events)

	return nil
}
