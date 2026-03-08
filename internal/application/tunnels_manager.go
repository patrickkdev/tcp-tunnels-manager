package application

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/patrickkdev/tcptunnel/internal/infrastructure/db"
	"github.com/patrickkdev/tcptunnel/internal/infrastructure/tcptunnels"
)

type Manager struct {
	tunnels        sync.Map // map[int]*tcptunnels.Tunnel
	tunnelRowsRepo db.TunnelRowsRepo
}

func NewManager(repo db.TunnelRowsRepo) *Manager {
	return &Manager{
		tunnelRowsRepo: repo,
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

			m.Reconcile(rows)
		}
	}
}

func (m *Manager) Reconcile(rows []db.TunnelRow) {
	// Start or update tunnels
	for port, row := range rows {

		listen := fmt.Sprintf(":%d", row.ListenPort)
		target := fmt.Sprintf("%s:%d", row.TargetHost, row.TargetPort)

		val, exists := m.tunnels.Load(port)

		if !exists {

			t := tcptunnels.New(listen, target)

			if err := t.Start(); err != nil {
				log.Printf("tunnel start failed port=%d err=%v", port, err)
				continue
			}

			m.tunnels.Store(port, t)

			log.Printf("tunnel started %s -> %s", listen, target)

			continue
		}

		t := val.(*tcptunnels.Tunnel)

		// restart if target changed
		if t.TargetAddr != target {

			log.Printf("tunnel target changed restarting port=%d", port)

			t.Stop()

			newTunnel := tcptunnels.New(listen, target)

			if err := newTunnel.Start(); err != nil {
				log.Printf("tunnel restart failed port=%d err=%v", port, err)
				continue
			}

			m.tunnels.Store(port, newTunnel)
		}
	}

	// Stop removed tunnels
	m.tunnels.Range(func(key, value any) bool {

		port := key.(int)
		tunnel := value.(*tcptunnels.Tunnel)

		if row, ok := rows[port]; !ok || !row.Enabled {

			log.Printf("tunnel stopping port=%d", port)

			tunnel.Stop()

			m.tunnels.Delete(port)
		}

		return true
	})
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