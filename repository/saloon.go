package repository

import (
	"database/sql"
	"seasaloon-backend-go/structs"
)
type SaloonRepository struct {
	DB *sql.DB
}

func GetAllSaloonCustomers(db *sql.DB, limit, offset int) (result []structs.Saloon, err error) {
	sql := `
        SELECT 
			s.id, s.name, s.location, s.open, s.close
		FROM saloon s
		WHERE s.is_delete = false
		LIMIT $1 OFFSET $2
    `

	rows, err := db.Query(sql, limit, offset)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var saloon structs.Saloon

		err = rows.Scan(
			&saloon.ID, &saloon.Name, &saloon.Location , &saloon.Open, &saloon.Close,
		)
		if err != nil {
			return
		}

		result = append(result, saloon)
	}

	return
}
func GetAllSaloon(db *sql.DB, limit, offset int) (result []structs.Saloon, err error) {
	sql := `
        SELECT 
			s.id, s.name, s.location, s.open, s.close, s.is_delete
		FROM saloon s
		LIMIT $1 OFFSET $2
    `

	rows, err := db.Query(sql, limit, offset)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var saloon structs.Saloon

		err = rows.Scan(
			&saloon.ID, &saloon.Name, &saloon.Location , &saloon.Open, &saloon.Close, &saloon.IsDelete,
		)
		if err != nil {
			return
		}

		result = append(result, saloon)
	}

	return
}

func InsertSaloon(db *sql.DB, saloon structs.Saloon) (err error) {
	sql := `
        INSERT INTO saloon (
            id, name, location, open, close, is_delete 
        ) VALUES ($1, $2, $3, $4, $5, $6)
    `

	_, err = db.Exec(sql,
		saloon.ID, saloon.Name, saloon.Location, saloon.Open, saloon.Close, false,
	)
	return
}


func UpdateSaloon(db *sql.DB, saloon structs.Saloon) (err error) {
    sql := `
        UPDATE saloon
        SET 
            name = $1, 
            location = $2, 
            open = $3, 
            close = $4
        WHERE id = $5
    `
    _, err = db.Exec(sql,
        saloon.Name, saloon.Location, saloon.Open, saloon.Close, saloon.ID,
    )
    return
}


func DeleteSaloon(db *sql.DB, saloon structs.Saloon) error {
    sql := `
        UPDATE saloon
        SET is_delete = true
        WHERE id = $1
    `
    _, err := db.Exec(sql, saloon.ID)
    if err != nil {
       panic(err)
    }

    return nil
}

func GetSaloonById(db *sql.DB, saloonID int) (saloon structs.Saloon, err error) {
    sql := `
        SELECT 
            id, name, location, open, close, is_delete
        FROM saloon
        WHERE id = $1
    `
    err = db.QueryRow(sql, saloonID).Scan(
        &saloon.ID, &saloon.Name, &saloon.Location, 
        &saloon.Open, &saloon.Close, &saloon.IsDelete,
    )
    
    return
}
