package tcptunnels

import (
	"context"
	"io"
	"net"
	"sync"
)

type Tunnel struct {
	ListenAddr string
	TargetAddr string

	listener net.Listener
	ctx      context.Context
	cancel   context.CancelFunc

	wg sync.WaitGroup
}

func New(listen, target string) *Tunnel {
	ctx, cancel := context.WithCancel(context.Background())

	return &Tunnel{
		ListenAddr: listen,
		TargetAddr: target,
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (t *Tunnel) Start() error {
	l, err := net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	t.listener = l

	t.wg.Add(1)
	go t.acceptLoop()

	return nil
}

func (t *Tunnel) acceptLoop() {
	defer t.wg.Done()

	for {
		conn, err := t.listener.Accept()
		if err != nil {
			select {
			case <-t.ctx.Done():
				return
			default:
				time.Sleep(50 * time.Millisecond)
				continue
			}
		}

		t.wg.Add(1)
		go t.handle(conn)
	}
}

func (t *Tunnel) handle(src net.Conn) {
	defer t.wg.Done()
	defer src.Close()

	dialer := &net.Dialer{
		Timeout: 10 * time.Second,
	}

	src, err := dialer.DialContext(t.ctx, "tcp", t.TargetAddr)
	if err != nil {
		return
	}
	defer src.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		io.Copy(dst, src)
		if tcp, ok := dst.(*net.TCPConn); ok {
			tcp.CloseWrite()
		}
		wg.Done()
	}()

	go func() {
		io.Copy(src, dst)
		if tcp, ok := src.(*net.TCPConn); ok {
			tcp.CloseWrite()
		}
		wg.Done()
	}()

	wg.Wait()
}

func (t *Tunnel) Stop() {
	t.cancel()

	if t.listener != nil {
		t.listener.Close()
	}

	t.wg.Wait()
}