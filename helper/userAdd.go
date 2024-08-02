package helper

import (
	"log"
	"userApiTest/database"
	"userApiTest/model"
)

func UserAdd() {
	database.Db.Exec("DELETE FROM users")

	test := model.User{
		Name:     "Correct",
		Email:    "correct@gmail.com",
		Password: "correct123",
	}
	if err := database.Db.Create(&test).Error; err != nil {
		log.Fatal(err.Error())
	}

}
