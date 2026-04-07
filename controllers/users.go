package controllers

import (
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
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			helpers.Error(c, http.StatusBadRequest, "invalid request")
			return
		}

		token, err := repository.RegisterUser(db, req.Username, req.Password)
		if err != nil {
			helpers.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		activationLink := "http://localhost:8080/activate?token=" + token

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

		err := repository.RegisterAdmin(db, req.Username, req.Password)
		if err != nil {
			helpers.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

	helpers.Success(c, http.StatusCreated, "Admin registered successfully", gin.H{})	
	}
}

func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			helpers.Error(c, http.StatusBadRequest, "invalid request")
			return
		}

		user, token, err := repository.LoginUser(db, req.Username, req.Password)
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

		err := db.QueryRow(`
			SELECT user_id FROM user_activation
			WHERE token = $1 AND expired_at > NOW()
		`, token).Scan(&userID)

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

		helpers.Success[any](c, 200, "account activated", nil)
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