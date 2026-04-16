package controllers

import (
	"bytes"
	"database/sql"
	"net/http"
	connection "seasaloon-backend-go/database/connections"
	"seasaloon-backend-go/helpers"
	"seasaloon-backend-go/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			helpers.Error(c, http.StatusBadRequest, "invalid request")
			return
		}

		token, err := repository.RegisterUser(db, req.Email, req.Password)
		if err != nil {
			helpers.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		activationLink := "http://localhost:8080/api/activate?token=" + token
		helpers.Success(c, http.StatusCreated, "User registered successfully", gin.H{
			"activation_link": activationLink,
		})
	}
}
func RegisterAdmin(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			helpers.Error(c, http.StatusBadRequest, "invalid request")
			return
		}

		err := repository.RegisterAdmin(db, req.Email, req.Password)
		if err != nil {
			helpers.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

	helpers.Success(c, http.StatusCreated, "Admin registered successfully", gin.H{})	
	}
}

func RegisterDoctor(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			helpers.Error(c, http.StatusBadRequest, "invalid request")
			return
		}

		err := repository.RegisterDoctor(db, req.Email, req.Password)
		if err != nil {
			helpers.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

	helpers.Success(c, http.StatusCreated, "Doctor registered successfully", gin.H{})	
	}
}

func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			helpers.Error(c, http.StatusBadRequest, "invalid request")
			return
		}

		user, token, err := repository.LoginUser(db, req.Email, req.Password)
		if err != nil {
			helpers.Error(c, http.StatusUnauthorized, err.Error())
			return
		}

		helpers.Success(c, http.StatusOK, "Login successfully", helpers.LoginResponse{
			Token: token,
			User:  user,
		})
	}
}
func ActivateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")

		var userID uuid.UUID
		var email string
		
		err := db.QueryRow(`
			SELECT ua.user_id, u.email
			FROM user_activation ua
			JOIN users u ON ua.user_id = u.id
			WHERE ua.token = $1 AND ua.expired_at > NOW()
		`, token).Scan(&userID, &email)

		if err != nil {
			helpers.Error(c, 400, "invalid or expired token")
			return
		}

		_, err = db.Exec(`
			UPDATE users SET is_active = true WHERE id = $1
		`, userID)

		if err != nil {
			helpers.Error(c, 500, "failed to activate user")
			return
		}
		_, _ = db.Exec(`
			DELETE FROM user_activation WHERE user_id = $1
		`, userID)
		go helpers.SendVerifiedEmail(email)
		c.File("helpers/email/verified.html")
	}
}

func SetCustomerMembership(c *gin.Context) {
	id := c.Param("id")

	userID, err := uuid.Parse(id)
	if err != nil {
		helpers.Error(c, 400, "invalid user id")
		return
	}

	err = repository.SetCustomerMembership(connection.DBConnections, userID)
	if err != nil {
		helpers.Error(c, 500, "failed to set customer membership")
		return
	}

	helpers.Success[any](c, 200, "customer membership updated", nil)
}

func GetAllCustomer(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	customers, err := repository.GetAllCustomers(connection.DBConnections, limit, offset)
	if err != nil {
		helpers.Error(c, 500, "failed to get customers")
		return
	}

	helpers.Success(c, http.StatusOK, "success get customers", customers)
}

func GetAllDoctor(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	doctors, err := repository.GetAllDoctors(connection.DBConnections, limit, offset)
	if err != nil {
		helpers.Error(c, 500, "failed to get doctors")
		return
	}

	helpers.Success(c, http.StatusOK, "success get doctors", doctors)
}

func SetProfile(c *gin.Context) {
	id := c.Param("id")

	userID, err := uuid.Parse(id)
	if err != nil {
		helpers.Error(c, 400, "invalid user id")
		return
	}
	file, _, err := c.Request.FormFile("avatar")
		if err != nil {
			helpers.Error(c, http.StatusBadRequest, "avatar is required")
			return
		}
		defer file.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(file)
		avatarBytes := buf.Bytes()
	err = repository.UpdateProfile(connection.DBConnections, userID, avatarBytes)
	if err != nil {
		helpers.Error(c, 500, "failed to update profile")
		return
	}

	helpers.Success[any](c, 200, "profile updated", nil)
}
func GetMe(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("user_id")
        if !exists {
            helpers.Error(c, http.StatusUnauthorized, "unauthorized")
            return
        }

        user, err := repository.GetUserByID(db, userID.(uuid.UUID))
        if err != nil {
            helpers.Error(c, http.StatusInternalServerError, err.Error())
            return
        }
        helpers.Success(c, http.StatusOK, "success", user)
    }
}