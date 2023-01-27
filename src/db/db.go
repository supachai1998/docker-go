package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	structsDB "docker_test/structs/db"
)

var Con *gorm.DB

func SetupDB() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok ",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))
	fmt.Println("DB", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	fmt.Println("con")
	if err != nil {
		return err
	}
	fmt.Println("con success")
	sql, _ := db.DB()
	sql.SetMaxIdleConns(25)
	sql.SetMaxOpenConns(50)

	// dropTable(db)
	AutoMigration(db)

	Con = db.Debug().Session(&gorm.Session{
		NewDB:       true,
		PrepareStmt: true,
	})

	return nil
}
func dropTable(db *gorm.DB) {
	db.Migrator().DropTable(&structsDB.Users{})
}
func AutoMigration(db *gorm.DB) {
	db.AutoMigrate(&structsDB.Users{})
}
