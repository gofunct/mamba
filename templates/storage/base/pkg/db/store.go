package db

import (
	"context"

	"{{[ .Project.Project ]}}/pkg/db/provider"
)

// Store design database interface with providers
type Store interface {
	Check() error
	Shutdown(ctx context.Context) error
	{{[- if .Project.Contract ]}}
	EventsProvider() provider.Events
	{{[- end ]}}
}
