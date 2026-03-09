package db

import (
	"context"

	"github.com/patrickkdev/tcptunnel/internal/domain"
)

type TunnelRowsRepo struct {
	dbConn *Connection
}

func NewTunnelRowsRepo(dbConn *Connection) *TunnelRowsRepo {
	return &TunnelRowsRepo{
		dbConn: dbConn,
	}
}

func (r *TunnelRowsRepo) List(ctx context.Context) ([]domain.TunnelRow, error) {
	var rows []domain.TunnelRow
	err := r.dbConn.SelectContext(ctx, &rows, "SELECT id, listen_port, target_host, target_port, enabled, created_at, updated_at FROM tcp_tunnels")
	if err != nil {
		return nil, err
	}

	return rows, nil
}
