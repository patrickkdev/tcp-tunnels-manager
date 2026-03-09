package application

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/patrickkdev/tcptunnel/internal/domain"
	"github.com/patrickkdev/tcptunnel/internal/infrastructure/db"
	"github.com/patrickkdev/tcptunnel/internal/infrastructure/tcptunnels"
)

type Manager struct {
	tunnels        sync.Map
	tunnelRowsRepo *db.TunnelRowsRepo
	tunnelLogsRepo *db.TunnelLogsRepo
}

func NewManager(tunnelRowsRepo *db.TunnelRowsRepo, tunnelLogsRepo *db.TunnelLogsRepo) *Manager {
	return &Manager{
		tunnelRowsRepo: tunnelRowsRepo,
		tunnelLogsRepo: tunnelLogsRepo,
	}
}

func (m *Manager) Run(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {

		case <-ctx.Done():
			log.Println("tunnel manager shutting down")

			m.shutdown()

			return

		case <-ticker.C:

			rows, err := m.tunnelRowsRepo.List(ctx)
			if err != nil {
				log.Println("failed to fetch tunnels:", err)
				continue
			}

			m.reconcile(ctx, rows)
		}
	}
}

func (m *Manager) reconcile(ctx context.Context, rows []domain.TunnelRow) {
	desired := buildDesired(rows)

	// CREATE or UPDATE
	for port, row := range desired {
		val, exists := m.tunnels.Load(port)

		// CREATE
		if !exists {
			t := tcptunnels.New(row)

			if err := t.Start(); err != nil {
				msg := fmt.Sprintf("tunnel start failed port=%d err=%v", port, err)
				m.logTunnel(ctx, row.ID, domain.LogLevelError, msg)
				continue
			}

			go m.consumeTunnelEvents(ctx, t)

			m.tunnels.Store(port, t)

			msg := fmt.Sprintf(
				"tunnel started %d -> %s:%d",
				port,
				row.TargetHost,
				row.TargetPort,
			)

			m.logTunnel(ctx, row.ID, domain.LogLevelInfo, msg)
			continue
		}

		// UPDATE
		t := val.(*tcptunnels.Tunnel)

		if t.TargetHost != row.TargetHost || t.TargetPort != row.TargetPort {

			msg := fmt.Sprintf("tunnel restarting due to config change port=%d", port)
			m.logTunnel(ctx, row.ID, domain.LogLevelWarning, msg)

			t.Stop()

			newTunnel := tcptunnels.New(row)

			if err := newTunnel.Start(); err != nil {
				msg := fmt.Sprintf("tunnel restart failed port=%d err=%v", port, err)
				m.logTunnel(ctx, row.ID, domain.LogLevelError, msg)
				continue
			}

			go m.consumeTunnelEvents(ctx, newTunnel)

			m.tunnels.Store(port, newTunnel)
		}
	}

	// DELETE
	m.tunnels.Range(func(key, value any) bool {
		port := key.(int)
		tunnel := value.(*tcptunnels.Tunnel)

		if _, exists := desired[port]; !exists {
			msg := fmt.Sprintf("tunnel stopped port=%d", port)
			m.logTunnel(ctx, tunnel.ID, domain.LogLevelInfo, msg)

			tunnel.Stop()

			m.tunnels.Delete(port)
		}

		return true
	})
}

func buildDesired(rows []domain.TunnelRow) map[int]domain.TunnelRow {
	desired := make(map[int]domain.TunnelRow)

	for _, r := range rows {
		if r.Enabled {
			desired[r.ListenPort] = r
		}
	}

	return desired
}

func (m *Manager) consumeTunnelEvents(ctx context.Context, t *tcptunnels.Tunnel) {
	for {
		select {

		case <-ctx.Done():
			return

		case e, ok := <-t.Events:
			if !ok {
				return
			}

			m.logTunnel(ctx, t.ID, e.Level, fmt.Sprintf("tunnel event: %s", e.Message))
		}
	}
}

func (m *Manager) logTunnel(ctx context.Context, tunnelID int, level domain.LogLevel, message string) {
	log.Printf("%d - %s: %s", tunnelID, level, message)

	logEntry, err := domain.NewTunnelLog(tunnelID, level, message)
	if err != nil {
		log.Println("failed creating tunnel log:", err)
		return
	}

	if err := m.tunnelLogsRepo.AddLog(ctx, logEntry); err != nil {
		log.Println("failed storing tunnel log:", err)
	}
}

func (m *Manager) shutdown() {
	m.tunnels.Range(func(key, value any) bool {
		port := key.(int)
		tunnel := value.(*tcptunnels.Tunnel)

		log.Printf("stopping tunnel port=%d", port)

		tunnel.Stop()

		return true
	})
}
