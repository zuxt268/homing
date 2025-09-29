package repository

import (
	"os"
	"testing"

	"github.com/zuxt268/homing/internal/infrastructure/database"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	conn, cleanUp := database.NewTestContainerDBClient()
	db = conn
	code := m.Run()
	cleanUp()
	os.Exit(code)
}
