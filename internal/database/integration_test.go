package database

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/miclle/goblet/internal/entity"
)

// TestDatabaseIntegration requires DSNs pointing to disposable test databases.
func TestDatabaseIntegration(t *testing.T) {
	for _, driver := range []struct {
		name string
		env  string
	}{
		{name: "postgres", env: "GOBLET_TEST_POSTGRES_DSN"},
		{name: "mysql", env: "GOBLET_TEST_MYSQL_DSN"},
	} {
		t.Run(driver.name, func(t *testing.T) {
			dsn := os.Getenv(driver.env)
			if dsn == "" {
				t.Skipf("set %s to a disposable test database", driver.env)
			}
			ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
			defer cancel()
			db, err := Open(ctx, driver.name, dsn)
			if err != nil {
				t.Fatalf("open: %v", err)
			}
			sqlDB, err := db.DB()
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := sqlDB.Close(); err != nil {
					t.Errorf("close database: %v", err)
				}
			}()
			for range 2 {
				if err := Migrate(ctx, db); err != nil {
					t.Fatalf("migrate: %v", err)
				}
			}

			tx := db.WithContext(ctx).Begin()
			if tx.Error != nil {
				t.Fatal(tx.Error)
			}
			defer tx.Rollback()
			item := entity.Example{Title: "Dependency compatibility", Body: "数据库验证"}
			if err := tx.Create(&item).Error; err != nil {
				t.Fatalf("create: %v", err)
			}
			if item.ID == 0 || item.CreatedAt.IsZero() {
				t.Fatalf("missing generated fields: %+v", item)
			}
			if err := tx.Model(&item).Update("title", "Updated").Error; err != nil {
				t.Fatalf("update: %v", err)
			}
			var got entity.Example
			if err := tx.First(&got, item.ID).Error; err != nil {
				t.Fatalf("read: %v", err)
			}
			if got.Title != "Updated" || got.Body != item.Body {
				t.Fatalf("unexpected record: %+v", got)
			}
			if err := tx.Delete(&item).Error; err != nil {
				t.Fatalf("delete: %v", err)
			}
			if err := tx.First(&entity.Example{}, item.ID).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
				t.Fatalf("deleted record should be hidden, got %v", err)
			}
		})
	}
}
