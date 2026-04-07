package controllers

import (
	"database/sql"
	"net/http"
	connection "seasaloon-backend-go/database/connections"
	"seasaloon-backend-go/repository"
	"seasaloon-backend-go/structs"

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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := repository.RegisterUser(db, req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
	}
}
func RegisterAdmin(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest

		// Validasi input JSON
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Call repository untuk register user
		err := repository.RegisterAdmin(db, req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Admin registered successfully"})
	}
}
func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest

		// Validasi input JSON
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, err := repository.LoginUser(db, req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
func ActivateUser(u *gin.Context) {
    var user structs.Users
    id := u.Param("id")

    userID, err := uuid.Parse(id)
    err = repository.ActivateUser(connection.DBConnections, userID)
    if err != nil {
       panic(err)
    }

    u.JSON(http.StatusOK, user)
}


func SetCustomerMembership(c *gin.Context) {
    var customer structs.Customer
    id := c.Param("id")

    userID, err := uuid.Parse(id)
    err = repository.SetCustomerMembership(connection.DBConnections, userID)
    if err != nil {
       panic(err)
    }

    c.JSON(http.StatusOK, customer)
}

func GetAllCustomer(s *gin.Context) {
    var (
       result gin.H
    )

    saloon, err := repository.GetAllCustomers(connection.DBConnections)

    if err != nil {
       result = gin.H{
          "result": err.Error(),
       }
    } else {
       result = gin.H{
          "result": saloon,
       }
    }

    s.JSON(http.StatusOK, result)
}