package db

import (
	"context"

	"github.com/patrickkdev/tcptunnel/internal/domain"
)

type TunnelLogsRepo struct {
	dbConn *Connection
}

func NewTunnelLogsRepo(dbConn *Connection) *TunnelLogsRepo {
	return &TunnelLogsRepo{
		dbConn: dbConn,
	}
}

func (r *TunnelLogsRepo) AddLog(ctx context.Context, log domain.TunnelLog) error {
	_, err := r.dbConn.ExecContext(ctx,
		"INSERT INTO tcp_tunnel_logs (tunnel_id, level, message) VALUES (?, ?, ?)",
		log.TunnelID, log.Level, log.Message,
	)

	return err
}
