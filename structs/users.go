package structs

import "github.com/google/uuid"

type Users struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Role     string    `json:"role"`
	IsActive bool      `json:"is_active"`
}

type Customer struct {
	CustomerID uuid.UUID `json:"customer_id"`
	UserID		uuid.UUID `json:"user_id"`
	User	   *Users	 `json:"user"`
	Status     string    `json:"status"`
}

type Admin struct {
	AdminID uuid.UUID `json:"admin_id"`
	UserID		uuid.UUID `json:"user_id"`
	User	   *Users	 `json:"user"`
}

type Doctor struct {
	DoctorID uuid.UUID `json:"doctor_id"`
	UserID		uuid.UUID `json:"user_id"`
	User	   *Users	 `json:"user"`
}