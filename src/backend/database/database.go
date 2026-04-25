package database

import (
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"github.com/oversteplab/oversteplab/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(dbPath string) *gorm.DB {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Enable WAL mode for better concurrent read performance
	DB.Exec("PRAGMA journal_mode=WAL")
	DB.Exec("PRAGMA foreign_keys=ON")

	return DB
}

func Close() {
	if DB != nil {
		sqlDB, _ := DB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}
}

func ResetDatabase(dbPath string) error {
	Close()
	if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// ReseedDatabase完全重置数据库：删除文件→重新初始化→迁移→填充种子数据
func ReseedDatabase(dbPath string) error {
	if err := ResetDatabase(dbPath); err != nil {
		return err
	}
	Init(dbPath)
	if err := AutoMigrate(); err != nil {
		return err
	}
	// 重新填充种子数据（seed包需导出Seed函数）
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

func AutoMigrate() error {
	return DB.AutoMigrate(
		&model.Company{},
		&model.User{},
		&model.VPSInstance{},
		&model.Order{},
		&model.Bill{},
		&model.Ticket{},
		&model.TicketReply{},
		&model.APIKey{},
		&model.AuditLog{},
		&model.Announcement{},
		&model.SystemConfig{},
	)
}
