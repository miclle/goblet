// Package service provides business logic and database operations.
package service

import (
	"context"
	"fmt"

	"github.com/fox-gonic/fox/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/miclle/go-single-page-binary-app/internal/entity"
	"github.com/miclle/go-single-page-binary-app/pkg/gormlog"
)

// Service holds the database connection and provides business logic methods.
type Service struct {
	db *gorm.DB
}

// New creates a new Service instance with the given driver and DSN.
// Supported drivers: "postgres", "mysql".
func New(ctx context.Context, driver, dsn string) (*Service, error) {
	l := logger.NewWithContext(ctx)

	var dialector gorm.Dialector
	switch driver {
	case "mysql":
		dialector = mysql.Open(dsn)
	default:
		dialector = postgres.Open(dsn)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   gormlog.New(0),
	})
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if err := db.WithContext(ctx).AutoMigrate(&entity.Example{}); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	l.Info("[Service] initialized")

	return &Service{db: db}, nil
}

// DB returns the underlying GORM database connection.
func (s *Service) DB() *gorm.DB {
	return s.db
}
