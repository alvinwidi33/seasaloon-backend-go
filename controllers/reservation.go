package controllers

import (
	"net/http"
	connection "seasaloon-backend-go/database/connections"
	"seasaloon-backend-go/repository"
	"seasaloon-backend-go/structs"
    "seasaloon-backend-go/helpers"
	"strconv"
	"github.com/gin-gonic/gin"
)

func GetAllReservation(r *gin.Context) {
	limitStr := r.DefaultQuery("limit", "10")
	offsetStr := r.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	data, err := repository.GetAllReservation(connection.DBConnections, limit, offset)

	if err != nil {
		helpers.Error(r, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.Success(r, http.StatusOK, "Success get all reservations", data)

}
func InsertReservation(r *gin.Context) {
    var reservation structs.Reservation

    err := r.BindJSON(&reservation)
    if err != nil {
		helpers.Error(r, http.StatusInternalServerError, err.Error())
		return
    }

    err = repository.InsertReservation(connection.DBConnections, reservation)
    if err != nil {
       	helpers.Error(r, http.StatusInternalServerError, err.Error())
		return
    }

    helpers.Success(r, http.StatusOK, "Reservation created successfully", gin.H{})
}

func CancelReservation(r *gin.Context) {
   var reservation structs.Reservation

   id, err := strconv.Atoi(r.Param("id"))
   if err != nil || id < 0 { 
	   helpers.Error(r, http.StatusBadRequest, "Invalid reservation ID")
       return
   }

   reservation.ID = id

   err = repository.CancelReservation(connection.DBConnections, reservation)
   if err != nil {
	   helpers.Error(r, http.StatusInternalServerError, err.Error())
       return
   }

   helpers.Success(r, http.StatusOK, "Reservation cancel successfully", gin.H{})
}

func DoneReservation(r *gin.Context) {
    var reservation structs.Reservation


    id, err := strconv.Atoi(r.Param("id"))
    if err != nil || id < 0 {
		helpers.Error(r, http.StatusBadRequest, "Invalid reservation ID")
        return
    }

    err = r.BindJSON(&reservation)
    if err != nil {
        helpers.Error(r, http.StatusBadRequest, "Invalid JSON Format")
        return
    }
    reservation.ID = id

    err = repository.DoneReservation(connection.DBConnections, reservation)
    if err != nil {
		helpers.Error(r, http.StatusBadRequest, err.Error())
        return
    }
	helpers.Success(r, http.StatusOK, "Reservation done successfully", gin.H{})
}


func GetAllReservationByCustomerID(r *gin.Context) {
	limitStr := r.DefaultQuery("limit", "10")
	offsetStr := r.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	customerID := r.Param("id")

	data, err := repository.GetAllReservationByCustomerID(
		connection.DBConnections,
		customerID,
		limit,
		offset,
	)

	if err != nil {
		helpers.Error(r, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.Success(r, http.StatusOK, "Success get reservations by customer", data)
}