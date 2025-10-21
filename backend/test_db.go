package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Category struct {
	ID           int
	CategoryID   int
	Name         string
	CategoryType int
}

func (Category) TableName() string {
	return "Category"
}

func main() {
	dsn := "root:password@tcp(db:3306)/mawinter?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Check which database we're connected to
	var dbName string
	db.Raw("SELECT DATABASE()").Scan(&dbName)
	fmt.Printf("Connected to database: %s\n", dbName)

	// Show tables
	var tables []string
	db.Raw("SHOW TABLES").Scan(&tables)
	fmt.Printf("Tables in database: %v\n", tables)

	var categories []Category
	result := db.Find(&categories)
	if result.Error != nil {
		log.Fatalf("Query failed: %v", result.Error)
	}

	fmt.Printf("Found %d categories\n", len(categories))
	for _, cat := range categories {
		fmt.Printf("- ID:%d CategoryID:%d Name:%s Type:%d\n", cat.ID, cat.CategoryID, cat.Name, cat.CategoryType)
	}
}
