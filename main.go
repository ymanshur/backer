package main

import (
	"backer/handler"
	"backer/user"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// refer https://github. .com/go-sal-driver/mysql#dsn-data-source-name for details
	dsn := "root:@tcp(127.0.0.1:3306)/backer?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)

	// userByEmail, err := userRepository.FindByEmail("ymanshur@gmail.com")
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }
	// if userByEmail.ID == 0 {
	// 	log.Fatalln("User tidak ditemukan")
	// } else {
	// 	fmt.Println(userByEmail.Name)
	// }

	userService := user.NewService(userRepository)

	input := user.LoginInput{
		Email:    "ymanshur@gmail.com",
		Password: "password",
	}
	user, err := userService.Login(input)
	if err != nil {
		fmt.Println("Terjadi kesalahan")
		log.Fatalln(err.Error())
	}
	fmt.Println(user.Email)
	fmt.Println(user.Name)

	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)

	router.Run()
}
