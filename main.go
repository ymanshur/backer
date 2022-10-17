package main

import (
	"backer/user"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// refer https://github. .com/go-sal-driver/mysql#dsn-data-source-name for details
	dsn := "root:@tcp(127.0.0.1:3306)/backer?charset=utf8mb4&parseTime-True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	// fmt.Println("Connection to database is good")

	userRepository := user.NewRepository(db)
	user := user.User{
		Name: "Test simpan",
	}

	userRepository.Save(user)

	// var users []user.User
	// db.Find(&users)

	// for _, user := range users {
	// 	fmt.Println(user.Name)
	// 	fmt.Println(user.Email)
	// 	fmt.Println("===============")
	// }

	// router := gin.Default()
	// router.GET("/", func(ctx *gin.Context) {
	// 	var users []user.User
	// 	db.Find(&users)

	// 	ctx.JSON(http.StatusOK, users)
	// })
	// router.Run()

	/**
	* Register endpoint:
	* 1. input
	* 2. handler: mapping input to struct
	* 3. service: mapping struct to User type
	* 4. repository: save struct User to database
	 */
}
