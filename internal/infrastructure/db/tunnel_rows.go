package db

type TunnelRowsRepo struct {
	db *DB
}

func NewTunnelRowsRepo(db *DB) *TunnelRowsRepo {
	return &TunnelRowsRepo{
		db: db,
	}
}

func (r *TunnelRowsRepo) List(ctx context.Context) ([]TunnelRow, error) {
	return []TunnelRow{}, nil
}