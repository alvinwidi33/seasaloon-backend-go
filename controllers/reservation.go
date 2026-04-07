package controllers

import (
	"log"
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

    r.JSON(http.StatusOK, reservation)
}

func CancelReservation(r *gin.Context) {
   var reservation structs.Reservation

   id, err := strconv.Atoi(r.Param("id"))
   if err != nil || id < 0 { 
       log.Println("Invalid reservation ID:", id, "Error:", err)
       r.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation ID"})
       return
   }

   reservation.ID = id

   err = repository.CancelReservation(connection.DBConnections, reservation)
   if err != nil {
       log.Println("Error updating reservation:", err)
       r.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }

   log.Println("Successfully cancel as done:", reservation)
   r.JSON(http.StatusOK, reservation)
}

func DoneReservation(r *gin.Context) {
    var reservation structs.Reservation


    id, err := strconv.Atoi(r.Param("id"))
    if err != nil || id < 0 {
        log.Println("Invalid reservation ID:", id, "Error:", err)
        r.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation ID"})
        return
    }

    err = r.BindJSON(&reservation)
    if err != nil {
        log.Println("JSON Binding Error:", err)
        r.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
        return
    }
    reservation.ID = id

    err = repository.DoneReservation(connection.DBConnections, reservation)
    if err != nil {
        log.Println("Error updating reservation:", err)
        r.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    log.Println("Successfully marked reservation as done, ID:", id)
    r.JSON(http.StatusOK, gin.H{"message": "Reservation marked as done", "id": id, "feedback": reservation.Feedback})
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