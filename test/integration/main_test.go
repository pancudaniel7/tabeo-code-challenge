package integration

import (
	"fmt"
	"os"
	"tabeo.org/challenge/internal/core/entity"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	dsn := os.Getenv("TABEO_MYSQL_TEST_DSN")
	if dsn == "" {
		dsn = "tabeo_user:password@tcp(127.0.0.1:3307)/tabeo?parseTime=true"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "connect failed:", err)
		os.Exit(1)
	}
	if err := db.AutoMigrate(&entity.Appointment{}); err != nil {
		fmt.Fprintln(os.Stderr, "migrate failed:", err)
		os.Exit(1)
	}
	if err := clean(db); err != nil {
		fmt.Fprintln(os.Stderr, "pre-clean failed:", err)
		os.Exit(1)
	}

	testDB = db
	code := m.Run()
	_ = clean(db)
	os.Exit(code)
}

func clean(db *gorm.DB) error {
	if err := db.Exec("TRUNCATE TABLE appointment").Error; err == nil {
		return nil
	}
	tx := db.Exec("SET FOREIGN_KEY_CHECKS=0").
		Exec("DELETE FROM appointment").
		Exec("SET FOREIGN_KEY_CHECKS=1")
	return tx.Error
}
