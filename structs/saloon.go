package structs
import "time"

type Saloon struct{
	ID 		 int 		`json:"id"`
	Name 	 string 	`json:"name"`
	Location string 	`json:"location"`
	Open 	 time.Time 	`json:"open"`
	Close 	 time.Time 	`json:"close"`
	IsDelete bool		`json:"is_delete"`
}