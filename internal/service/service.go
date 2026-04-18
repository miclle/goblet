// Package service provides business logic and database operations.
package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/fox-gonic/fox/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/miclle/go-single-page-binary-app/internal/entity"
	"github.com/miclle/go-single-page-binary-app/pkg/gormlog"
)

// Service holds the database connection and provides business logic methods.
type Service struct {
	db      *gorm.DB
	dialect string // "postgres" or "mysql"
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

	return &Service{
		db:      db,
		dialect: driver,
	}, nil
}

// DB returns the underlying GORM database connection.
func (s *Service) DB() *gorm.DB {
	return s.db
}

// isMySQL returns true if the current database is MySQL.
func (s *Service) isMySQL() bool {
	return s.dialect == "mysql"
}

// whereLike adds a case-insensitive LIKE condition across multiple columns (OR).
// PostgreSQL uses ILIKE; MySQL LIKE is case-insensitive by default with utf8mb4 collation.
func (s *Service) whereLike(query *gorm.DB, columns []string, pattern string) *gorm.DB {
	op := "ILIKE"
	if s.isMySQL() {
		op = "LIKE"
	}
	conds := make([]string, len(columns))
	args := make([]any, len(columns))
	for i, col := range columns {
		conds[i] = fmt.Sprintf("%s %s ?", col, op)
		args[i] = pattern
	}
	return query.Where(strings.Join(conds, " OR "), args...)
}
