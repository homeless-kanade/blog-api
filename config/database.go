package config

import (
	"blog-api/models"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Sử dụng SQLite lưu thành file local tên là blog.db
	database, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})

	if err != nil {
		panic("Không thể kết nối đến cơ sở dữ liệu!")
	}

	// Tự động tạo bảng dựa trên struct (Auto Migration)
	err = database.AutoMigrate(&models.User{}, &models.Post{})
	if err != nil {
		fmt.Println("Migration lỗi:", err)
	}

	DB = database
	fmt.Println("✅ Kết nối Database thành công và Auto Migrate hoàn tất!")
}
