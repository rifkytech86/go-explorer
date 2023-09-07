package domain

import (
	"context"
)

type MigrationUseCase interface {
	InitialDB(c context.Context) error
}

type MigrationRepository interface {
	InitialDB(c context.Context) error
}
