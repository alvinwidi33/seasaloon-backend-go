package main

import (
	"database/sql"
	"os"
	"seasaloon-backend-go/controllers"
	connection "seasaloon-backend-go/database/connections"
	migration "seasaloon-backend-go/database/migrations"
	"seasaloon-backend-go/middleware"

	"github.com/gin-gonic/gin"
)
var (
    DB *sql.DB
)



func main() {
	connection.Initiator()

	DB = connection.DBConnections
	if DB == nil {
		panic("Database connection is nil")
	}


	defer DB.Close()
	migration.Initiator(DB)

	router := gin.Default()

	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware()) 
	{
		protected.GET("/reservation", middleware.AuthMiddleware("Customer", "Admin", "Doctor"), controllers.GetAllReservation)
		protected.GET("/reservation/:id", middleware.AuthMiddleware("Customer", "Admin", "Doctor"), controllers.GetAllReservationByCustomerID)

		protected.POST("/reservation", middleware.AuthMiddleware("Customer"), controllers.InsertReservation)
		protected.PATCH("/reservation/:id/cancel", middleware.AuthMiddleware("Customer"), controllers.CancelReservation)
		protected.PATCH("/reservation/:id/done", middleware.AuthMiddleware("Customer"), controllers.DoneReservation)

		protected.GET("/saloon/customer", middleware.AuthMiddleware("Customer"), controllers.GetAllSaloonCustomers)
		protected.POST("/saloon", middleware.AuthMiddleware("Admin"), controllers.InsertSaloon)
		protected.PUT("/saloon/:id", middleware.AuthMiddleware("Admin"), controllers.UpdateSaloon)
		protected.PATCH("/saloon/:id/delete", middleware.AuthMiddleware("Admin"), controllers.DeleteSaloon)
        protected.GET("/saloon/:id", middleware.AuthMiddleware("Admin","Customer"), controllers.GetSaloonById)

		protected.GET("/users", middleware.AuthMiddleware("Admin"), controllers.GetAllCustomer)
		protected.PATCH("/users/:id/member", middleware.AuthMiddleware("Admin"), controllers.SetCustomerMembership)
        protected.POST("/admin", middleware.AuthMiddleware("Admin"), controllers.RegisterAdmin(DB))
        protected.GET("/admin/saloon", middleware.AuthMiddleware("Admin"), controllers.GetAllSaloon)
	}
	router.POST("/api/login", controllers.Login(DB))
	router.POST("/api/register", controllers.Register(DB))
	router.GET("/api/activate", controllers.ActivateUser(DB))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}