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
       helpers.Error(s, http.StatusInternalServerError, err.Error())
    }

    err = repository.InsertSaloon(connection.DBConnections, saloon)
    if err != nil {
       helpers.Error(s, http.StatusInternalServerError, err.Error())
    }

    helpers.Success(s, http.StatusCreated, "Saloon created successfully", gin.H{})	
}

func UpdateSaloon(s *gin.Context) {
   var saloon structs.Saloon
   id, _ := strconv.Atoi(s.Param("id"))

   err := s.BindJSON(&saloon)
   if err != nil {
       helpers.Error(s, http.StatusBadRequest, "Invalid saloon ID")
       return
   }

   saloon.ID = id

   err = repository.UpdateSaloon(connection.DBConnections, saloon)
   if err != nil {
       helpers.Error(s, http.StatusInternalServerError, err.Error())
       return
   }

   helpers.Success(s, http.StatusOK, "Saloon updated successfully", gin.H{})
}


func DeleteSaloon(s *gin.Context) {
    var saloon structs.Saloon
    id, _ := strconv.Atoi(s.Param("id"))

    saloon.ID = id
    err := repository.DeleteSaloon(connection.DBConnections, saloon)
    if err != nil {
       helpers.Error(s, http.StatusBadRequest, err.Error())
    }

    helpers.Success(s, http.StatusOK, "Saloon deleted successfully", gin.H{})
}

func GetSaloonById(s *gin.Context) {
	id, err := strconv.Atoi(s.Param("id"))
	if err != nil {
		helpers.Error(s, http.StatusBadRequest, "invalid saloon id")
		return
	}

	data, err := repository.GetSaloonById(connection.DBConnections, id)
	if err != nil {
		helpers.Error(s, http.StatusInternalServerError, "failed to get saloon")
		return
	}

	helpers.Success(s, http.StatusOK, "Success get saloon", data)
}