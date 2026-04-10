package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"seasaloon-backend-go/structs"
	"time"
)

type ReservationRepository struct {
	DB *sql.DB
}

func GetAllReservation(db *sql.DB, limit, offset int) (result []structs.Reservation, err error) {
	sqlQuery := `
        SELECT 
            r.id, r.services, r.start, r.done, 
            r.is_done, r.is_cancel, r.rating, r.feedback, 
            c.customer_id, c.status, u.id, u.email,
            u.role, u.is_active, s.id, s.name, s.location, s.open, s.close
        FROM reservation r
        JOIN customer c ON r.customer_id = c.customer_id
        JOIN users u ON c.user_id = u.id
		JOIN saloon s ON r.saloon_id = s.id
        LIMIT $1 OFFSET $2
    `

	rows, err := db.Query(sqlQuery, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()


	for rows.Next() {
		var reservation structs.Reservation
		var customer structs.Customer
		var user structs.Users
		var saloon structs.Saloon
		var servicesData []byte

		err = rows.Scan(
			&reservation.ID, &servicesData, &reservation.Start, &reservation.Done, &reservation.IsDone,
			&reservation.IsCancel, &reservation.Rating, &reservation.Feedback, &customer.CustomerID, &customer.Status, &user.ID,
			&user.Email, &user.Role, &user.IsActive, &saloon.ID, &saloon.Name, &saloon.Location, &saloon.Open, &saloon.Close,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}

		if err := json.Unmarshal(servicesData, &reservation.Services); err != nil {
			return nil, fmt.Errorf("failed to unmarshal services JSON: %w", err)
		}


		customer.User = &user
		reservation.Customer = &customer
		reservation.Saloon = &saloon
		result = append(result, reservation)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return result, nil
}

func GetAllReservationByCustomerID(db *sql.DB, customerID string, limit, offset int) (result []structs.Reservation, err error) {
	sqlQuery := `
        SELECT 
			r.id, r.services, r.start, r.done, 
			r.is_done, r.is_cancel, r.rating, r.feedback, 
			c.customer_id, c.status , u.id AS user_id, u.email,
			u.role, u.is_actives.id, s.name, s.location, s.open, s.close
		FROM reservation r
		JOIN customer c ON r.customer_id = c.customer_id
		JOIN users u ON c.user_id = u.id
		JOIN saloon s ON r.saloon_id = s.id
		WHERE c.customer_id = $1
        LIMIT $2 OFFSET $3
    `

	rows, err := db.Query(sqlQuery, customerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var reservation structs.Reservation
		var customer structs.Customer
		var user structs.Users
		var saloon structs.Saloon
		var servicesData []byte 

		err = rows.Scan(
			&reservation.ID, &servicesData, &reservation.Start, &reservation.Done, &reservation.IsDone,
			&reservation.IsCancel, &reservation.Rating, &reservation.Feedback, &customer.CustomerID, &customer.Status, &user.ID,
			&user.Email, &user.Role, &user.IsActive, &saloon.ID, &saloon.Name, &saloon.Location, &saloon.Open, &saloon.Close,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}

		if err := json.Unmarshal(servicesData, &reservation.Services); err != nil {
			return nil, fmt.Errorf("failed to unmarshal services JSON: %w", err)
		}

		customer.User = &user
		reservation.Customer = &customer
		reservation.Saloon = &saloon
		result = append(result, reservation)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return result, nil
}

func InsertReservation(db *sql.DB, reservation structs.Reservation) error {
    var isMember string
    err := db.QueryRow(`SELECT status FROM customer WHERE customer_id = $1`, reservation.CustomerID).Scan(&isMember)
    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("Customer not found with ID: %s", reservation.CustomerID)
        }
        return err
    }

    if isMember == "Not Member" {
        return fmt.Errorf("Reservation not allowed for non-members")
    }

    var openTime, closeTime time.Time
    var isDelete *bool
    err = db.QueryRow(`SELECT open, close, is_delete FROM saloon WHERE id = $1`, reservation.SaloonID).Scan(&openTime, &closeTime, &isDelete)
    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("saloon not found with ID: %d", reservation.SaloonID)
        }
        return err
    }

    if isDelete != nil && *isDelete {
        return fmt.Errorf("reservation not allowed, saloon is deleted")
    }

    localReservationStart := reservation.Start.UTC()
    localOpenTime := openTime.UTC()
    localCloseTime := closeTime.UTC()

    log.Printf("Reservation start time (UTC): %v, Saloon open: %v, close: %v\n", localReservationStart, localOpenTime, localCloseTime)
    reservationEndTime := localReservationStart.Add(time.Hour)


    servicesJSON, err := json.Marshal(reservation.Services)
    if err != nil {
        return err
    }

    sql := `
        INSERT INTO reservation (
            id, services, start, done, is_done, is_cancel, rating, feedback, customer_id, saloon_id
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
        )
    `
    _, err = db.Exec(sql, reservation.ID, servicesJSON, localReservationStart, reservationEndTime,
        reservation.IsDone, reservation.IsCancel, reservation.Rating, reservation.Feedback, reservation.CustomerID,
        reservation.SaloonID,
    )

    if err != nil {
        log.Printf("Error inserting reservation: %v\n", err)
        return err
    }

    log.Printf("Reservation created with ID: %d\n", reservation.ID)
    return nil
}

func CancelReservation(db *sql.DB, reservation structs.Reservation) error {
    sql := `
        UPDATE reservation
		SET is_cancel = true
		WHERE id = $1
    `
    _, err := db.Exec(sql, reservation.ID)
    if err != nil {
       panic(err)
    }

    return nil
}

func DoneReservation(db *sql.DB, reservation structs.Reservation) error {
    log.Println("Updating reservation with ID:", reservation.ID)

    sqlQuery := `
        UPDATE reservation
        SET is_done = true, rating = $2, feedback = $3
        WHERE id = $1
    `

    result, err := db.Exec(sqlQuery, reservation.ID, reservation.Rating, reservation.Feedback)
    if err != nil {
        log.Println("DB Error:", err)
        return fmt.Errorf("failed to update reservation: %w", err)
    }

    rowsAffected, _ := result.RowsAffected()
    log.Println("Rows affected:", rowsAffected)

    if rowsAffected == 0 {
        return fmt.Errorf("no reservation found with ID %d", reservation.ID)
    }

    return nil
}