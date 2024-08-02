package test

import (
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
	"github.com/stretchr/testify/require"
)

func TestUserShow(t *testing.T) {
	mock := database.DbSet()
	defer database.MockDB.Close()
	gin.SetMode(gin.TestMode)

	t.Run("Successfull fetch", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE \"users\".\"deleted_at\" IS NULL").WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email", "password"}).
				AddRow(1, "user1", "user1@gmail.com", "user1@123").
				AddRow(2, "user2", "user2@gmail.com", "user2@123"),
		)
		router := gin.Default()
		router.GET("/users", controllers.UserShow)

		req, _ := http.NewRequest("GET", "/users", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string][]model.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Len(t, response["users"], 2)
		assert.Equal(t, "user1", response["users"][0].Name)
		assert.Equal(t, "user1@gmail.com", response["users"][0].Email)
		assert.Equal(t, "user2", response["users"][1].Name)
		assert.Equal(t, "user2@gmail.com", response["users"][1].Email)

	})
}
