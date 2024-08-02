package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"userApiTest/controllers"
	"userApiTest/database"
	"userApiTest/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestEditUser(t *testing.T) {
	mock := database.DbSet()
	gin.SetMode(gin.TestMode)
	defer database.MockDB.Close()

	t.Run("Succesfull edit", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users" SET "email"=\$1,"name"=\$2,"password"=\$3,"updated_at"=\$4 WHERE id = \$5 AND "users"."deleted_at" IS NULL`).
			WithArgs("userEdit1@gmail.com", "userEdit1", sqlmock.AnyArg(), sqlmock.AnyArg(), 11).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		r := gin.Default()
		r.PATCH("/user", controllers.UserEdit)

		user := model.User{
			Name:     "userEdit1",
			Email:    "userEdit1@gmail.com",
			Password: "userEdit1@123",
		}
		jsonvalue, _ := json.Marshal(user)
		req, _ := http.NewRequest("PATCH", "/user?id=11", bytes.NewBuffer(jsonvalue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "User updated successfully")
	})

}
