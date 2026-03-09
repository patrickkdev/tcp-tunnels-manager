package db

import (
	"context"

	"github.com/patrickkdev/tcptunnelsmanager/internal/domain"
)

type TunnelRowsRepo struct {
	dbConn *Connection
}

func NewTunnelRowsRepo(dbConn *Connection) *TunnelRowsRepo {
	return &TunnelRowsRepo{
		dbConn: dbConn,
	}
}

type tunnelRow struct {
	ID         int    `db:"id"`
	ListenPort int    `db:"listen_port"`
	TargetHost string `db:"target_host"`
	TargetPort int    `db:"target_port"`
	Enabled    bool   `db:"enabled"`
	CreatedAt  string `db:"created_at"`
	UpdatedAt  string `db:"updated_at"`
}

func (r *TunnelRowsRepo) List(ctx context.Context) ([]domain.TunnelRow, error) {
	var rows []tunnelRow
	err := r.dbConn.SelectContext(ctx, &rows, "SELECT id, listen_port, target_host, target_port, enabled, created_at, updated_at FROM tcp_tunnels")
	if err != nil {
		return nil, err
	}

	return mapTunnelRows(rows), nil
}

func mapTunnelRows(rows []tunnelRow) []domain.TunnelRow {
	tunnelRows := make([]domain.TunnelRow, len(rows))

	for i, r := range rows {
		tunnelRows[i] = domain.TunnelRow{
			ID:         r.ID,
			ListenPort: r.ListenPort,
			TargetHost: r.TargetHost,
			TargetPort: r.TargetPort,
			Enabled:    r.Enabled,
			CreatedAt:  r.CreatedAt,
			UpdatedAt:  r.UpdatedAt,
		}
	}

	return tunnelRows
}
