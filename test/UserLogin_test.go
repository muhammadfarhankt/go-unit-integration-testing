// package test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"userApiTest/controllers"
// 	"userApiTest/database"
// 	"userApiTest/model"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func SetupTestDBSignin(t *testing.T) sqlmock.Sqlmock {
// 	mock := database.DbSet()
// 	mock.ExpectBegin()
// 	mock.ExpectExec("DELETE FROM users").WillReturnResult(sqlmock.NewResult(0, 0))
// 	database.Db.Exec("DELETE FROM users")
// 	testUser := model.User{Name: "Correct", Email: "correct@gmail.com", Password: "correct123"}
// 	if err := database.Db.Create(&testUser).Error; err != nil {
// 		t.Logf(err.Error())
// 	}

// 	return mock
// }

// func TestUserLogin(t *testing.T) {
// 	mock := SetupTestDBSignin(t)
// 	defer database.MockDB.Close()
// 	gin.SetMode(gin.TestMode)

// 	t.Run("Valid Login", func(t *testing.T) {

// 		mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE email=\\$1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT \\$2").
// 			WithArgs("correct@gmail.com", 1).
// 			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password"}).
// 				AddRow(1, "Correct", "correct@gmail.com", "correct123"))

// 		r := gin.Default()
// 		r.POST("/user", controllers.UserLogin)
// 		payLoad := model.User{Email: "user@gmail.com", Password: "user@123"}
// 		jsonValue, _ := json.Marshal(payLoad)

// 		req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonValue))
// 		if err != nil {
// 			t.Logf(err.Error())
// 		}

// 		req.Header.Set("Content-Type", "application/json")

// 		w := httptest.NewRecorder()
// 		r.ServeHTTP(w, req)

// 		if w.Code != http.StatusOK {
// 			t.Logf("Response Code: %d", w.Code)
// 			t.Logf("Response Body: %s", w.Body.String())
// 		}
// 		require.Equal(t, http.StatusOK, w.Code)

// 		var response map[string]string
// 		err = json.Unmarshal(w.Body.Bytes(), &response)
// 		require.NoError(t, err)
// 		assert.Equal(t, "Logged in successfully", response["message"])

//			if err := mock.ExpectationsWereMet(); err != nil {
//				t.Errorf("mock expectations were not met: %v", err)
//			}
//		})
//	}
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
	"golang.org/x/crypto/bcrypt"
)

func TestSingnIn(t *testing.T) {
	mock := database.DbSet()
	defer database.MockDB.Close()

	gin.SetMode(gin.TestMode)

	t.Run("successful login", func(t *testing.T) {
		password, _ := bcrypt.GenerateFromPassword([]byte("user@123"), bcrypt.DefaultCost)

		mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE email=\\$1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT \\$2").
			WithArgs("user@gmail.com", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password"}).
				AddRow(1, "user", "user@gmail.com", password))

		router := gin.Default()
		router.POST("/signin", controllers.UserLogin)
		loginInput := model.User{
			Email:    "user@gmail.com",
			Password: "user@123",
		}
		jsonValue, _ := json.Marshal(loginInput)
		req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Logged in successfully")
	})
	t.Run("wrong email", func(t *testing.T) {

		mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE email=\\$1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT \\$2").
			WithArgs("wrong@gmail.com", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password"}))

		router := gin.Default()
		router.POST("/signin", controllers.UserLogin)
		loginInput := model.User{
			Email:    "wrong@gmail.com",
			Password: "user@123",
		}
		jsonValue, _ := json.Marshal(loginInput)
		req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "User not found")
	})
	// t.Run("wrong password", func(t *testing.T) {
	// 	password, _ := bcrypt.GenerateFromPassword([]byte("user@123"), bcrypt.DefaultCost)
	// 	mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE email=\\$1 AND \"users\".\"deleted_at\" IS NULL ORDER BY \"users\".\"id\" LIMIT \\$2").
	// 		WithArgs("user@gmail.com", 1).
	// 		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password"}).
	// 			AddRow(1, "user", "user@gmail.com", password))

	// 	router := gin.Default()
	// 	router.POST("/signin", controllers.UserLogin)
	// 	loginInput := model.User{
	// 		Email:    "user@gmail.com",
	// 		Password: "wrong",
	// 	}
	// 	jsonValue, _ := json.Marshal(loginInput)
	// 	req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(jsonValue))
	// 	req.Header.Set("Content-Type", "application/json")

	// 	w := httptest.NewRecorder()
	// 	router.ServeHTTP(w, req)

	// 	assert.Equal(t, http.StatusBadRequest, w.Code)
	// 	assert.Contains(t, w.Body.String(), "Wrong email or password")
	// })
}
