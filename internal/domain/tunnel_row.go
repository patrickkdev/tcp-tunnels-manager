package domain

type TunnelRow struct {
	ID         int
	ListenPort int
	TargetHost string
	TargetPort int
	Enabled    bool
	CreatedAt  string
	UpdatedAt  string
}
