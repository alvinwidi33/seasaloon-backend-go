package repository

import (
	"database/sql"
	"errors"
	"os"
	"seasaloon-backend-go/structs"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(userID uuid.UUID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}


func ParseToken(tokenStr string) (uuid.UUID, string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	if err != nil {
		return uuid.UUID{}, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, err := uuid.Parse(claims["user_id"].(string))
		if err != nil {
			return uuid.UUID{}, "", err
		}

		role, ok := claims["role"].(string)
		if !ok {
			return uuid.UUID{}, "", errors.New("role not found in token")
		}

		return userID, role, nil
	}

	return uuid.UUID{}, "", jwt.ErrSignatureInvalid
}

func RegisterUser(db *sql.DB, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newID := uuid.New()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() 

	_, err = tx.Exec(`
		INSERT INTO users (id, username, password, role, is_active) 
		VALUES ($1, $2, $3, $4, $5)
	`, newID, username, string(hashedPassword), "Customer", true)
	if err != nil {
		return errors.New("failed to register user: " + err.Error())
	}

	_, err = tx.Exec(`
		INSERT INTO customer (customer_id, user_id, status) 
		VALUES ($1, $2, $3)
	`, uuid.New(), newID, "Not Member")
	if err != nil {
		return errors.New("failed to create customer: " + err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
func RegisterAdmin(db *sql.DB, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newID := uuid.New()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() 

	_, err = tx.Exec(`
		INSERT INTO users (id, username, password, role, is_active) 
		VALUES ($1, $2, $3, $4, $5)
	`, newID, username, string(hashedPassword), "Admin", true)
	if err != nil {
		return errors.New("failed to register user: " + err.Error())
	}

	_, err = tx.Exec(`
		INSERT INTO admin (admin_id, user_id) 
		VALUES ($1, $2)
	`, uuid.New(), newID)
	if err != nil {
		return errors.New("failed to create admin: " + err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
// Login User
func LoginUser(db *sql.DB, username, password string) (string, error) {
	var user structs.Users

	err := db.QueryRow("SELECT id, password, is_active, role FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Password, &user.IsActive, &user.Role)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	if !user.IsActive {
		return "", errors.New("user is no longer active")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	token, err := GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}


func ActivateUser(db *sql.DB, userID uuid.UUID) error {
	var currentStatus bool
	sql := `SELECT is_active FROM users WHERE id = $1`
	err := db.QueryRow(sql, userID).Scan(&currentStatus)
	if err != nil {
		return errors.New("failed to fetch activate: " + err.Error())
	}

	var newStatus bool
	if currentStatus == true {
		newStatus = false
	} else {
		newStatus = true
	}

	updateSQL := `
		UPDATE users
		SET is_active = $1
		WHERE id = $2
	`
	_, err = db.Exec(updateSQL, newStatus, userID)
	if err != nil {
		return errors.New("failed to update membership status: " + err.Error())
	}

	return nil
}

func SetCustomerMembership(db *sql.DB, userID uuid.UUID) error {
	var currentStatus string
	sql := `SELECT status FROM customer WHERE user_id = $1`
	err := db.QueryRow(sql, userID).Scan(&currentStatus)
	if err != nil {
		return errors.New("failed to fetch customer membership status: " + err.Error())
	}

	var newStatus string
	if currentStatus == "Member" {
		newStatus = "Not Member"
	} else {
		newStatus = "Member"
	}

	updateSQL := `
		UPDATE customer
		SET status = $1
		WHERE user_id = $2
	`
	_, err = db.Exec(updateSQL, newStatus, userID)
	if err != nil {
		return errors.New("failed to update membership status: " + err.Error())
	}

	return nil
}


func GetAllCustomers(db *sql.DB) (result []structs.Customer, err error) {
	sql := `
        SELECT  
			c.customer_id, c.status , u.id, u.username,
			u.role, u.is_active
		FROM customer c
		JOIN users u ON c.user_id = u.id
    `

	rows, err := db.Query(sql)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var customer structs.Customer
		var user structs.Users

		err = rows.Scan(
			&customer.CustomerID, &customer.Status, &user.ID,
			&user.Username, &user.Role, &user.IsActive,
		)
		if err != nil {
			return
		}
		customer.User = &user
		result = append(result, customer)
	}

	return
}