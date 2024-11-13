package main

import (
	"go/adv-demo/internal/link"
	"go/adv-demo/internal/user"
	"os"

	"github.com/lpernett/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil{
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil{
		panic(err.Error())
	}
	
	db.AutoMigrate(&link.Link{}, &user.User{})
}