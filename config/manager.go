package config

var global Manager

//go:generate mockgen -source=manager.go -destination=mock_manager.go -package=config
type Manager interface {
	ExternalURL() ExternalURL
	Server() Server
	Postgres() Postgres
}

type manager struct {
	config *config
}

func (m *manager) ExternalURL() ExternalURL {
	return m.config.ExternalURL
}

func (m *manager) Server() Server {
	return m.config.Server
}

func (m *manager) Postgres() Postgres {
	return m.config.Postgres
}
