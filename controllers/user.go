package controllers

import (
	"log"
	"net/http"
	"strconv"
	"userApiTest/database"
	"userApiTest/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func UserLogin(c *gin.Context) {

	var email model.User
	var user model.User

	c.ShouldBindJSON(&email)

	if err := database.Db.Where("email=?", email.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(email.Password)); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Wrong email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})

}

func UserShow(c *gin.Context) {

	var user []model.User

	if err := database.Db.Find(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": user,
	})
}

func UserEdit(c *gin.Context) {

	var email model.User

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	c.BindJSON(&email)
	Password, _ := bcrypt.GenerateFromPassword([]byte(email.Password), 8)
	email.Password = string(Password)

	updates := map[string]interface{}{
		"name":     email.Name,
		"email":    email.Email,
		"password": email.Password,
	}
	if err := database.Db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})

}

func UserSignup(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Failed to bind user:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	Password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	user.Password = string(Password)

	if err := database.Db.Create(&user).Error; err != nil {
		log.Println("Failed to create user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully now"})
}
