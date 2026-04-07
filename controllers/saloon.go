package controllers

import (
	"net/http"
	connection "seasaloon-backend-go/database/connections"
	"seasaloon-backend-go/repository"
	"seasaloon-backend-go/structs"
	"strconv"
   "seasaloon-backend-go/helpers"
	"github.com/gin-gonic/gin"
)

func GetAllSaloon(s *gin.Context) {
   limitStr := s.DefaultQuery("limit", "10")
	offsetStr := s.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

   data, err := repository.GetAllSaloon(connection.DBConnections, limit, offset)

	if err != nil {
		helpers.Error(s, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.Success(s, http.StatusOK, "Success get all Saloons", data)
}

func GetAllSaloonCustomers(s *gin.Context) {
   limitStr := s.DefaultQuery("limit", "10")
	offsetStr := s.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

   data, err := repository.GetAllSaloonCustomers(connection.DBConnections, limit, offset)

	if err != nil {
		helpers.Error(s, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.Success(s, http.StatusOK, "Success get all Saloons", data)
}

func InsertSaloon(s *gin.Context) {
    var saloon structs.Saloon

    err := s.BindJSON(&saloon)
    if err != nil {
       panic(err)
    }

    err = repository.InsertSaloon(connection.DBConnections, saloon)
    if err != nil {
       panic(err)
    }

    s.JSON(http.StatusOK, saloon)
}

func UpdateSaloon(s *gin.Context) {
   var saloon structs.Saloon
   id, _ := strconv.Atoi(s.Param("id"))

   err := s.BindJSON(&saloon)
   if err != nil {
       s.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
       return
   }

   saloon.ID = id


   err = repository.UpdateSaloon(connection.DBConnections, saloon)
   if err != nil {
       s.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }

   s.JSON(http.StatusOK, saloon)
}



func DeleteSaloon(s *gin.Context) {
    var saloon structs.Saloon
    id, _ := strconv.Atoi(s.Param("id"))

    saloon.ID = id
    err := repository.DeleteSaloon(connection.DBConnections, saloon)
    if err != nil {
       panic(err)
    }

    s.JSON(http.StatusOK, saloon)
}

func GetSaloonById(s *gin.Context) {
	id, err := strconv.Atoi(s.Param("id"))
	if err != nil {
		s.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	saloon, err := repository.GetSaloonById(connection.DBConnections, id)
	if err != nil {
		s.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	s.JSON(http.StatusOK, saloon)
}