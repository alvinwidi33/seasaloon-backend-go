package structs

import (
	"time"
	"github.com/google/uuid"
)

type Reservation struct {
	ID       int                    `json:"id"`
	Services map[string]interface{} `json:"services"` 
	Start    time.Time              `json:"start"`
	Done     time.Time              `json:"done"`
	IsDone   bool                   `json:"is_done"`
	IsCancel bool                   `json:"is_cancel"`
	Rating   int                    `json:"rating"`
	Feedback *string                 `json:"feedback"`
	CustomerID uuid.UUID			`json:"customer_id"`
	Customer *Customer				`json:"customer"`
	SaloonID int					`json:"saloon_id"`
	Saloon	 *Saloon				`json:"saloon"`
}