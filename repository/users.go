package repository

import (
	"database/sql"
	"errors"
	"os"
	"seasaloon-backend-go/structs"
	"time"
	"crypto/rand"
	"encoding/hex"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"seasaloon-backend-go/helpers"
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
func GenerateActivationToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	return hex.EncodeToString(bytes), err
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

func RegisterUser(db *sql.DB, email, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	newID := uuid.New()

	tx, err := db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO users (id, email, password, role, is_active) 
		VALUES ($1, $2, $3, $4, $5)
	`, newID, email, string(hashedPassword), "Customer", false)
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(`
		INSERT INTO customer (customer_id, user_id, status) 
		VALUES ($1, $2, $3)
	`, uuid.New(), newID, "Not Member")
	if err != nil {
		return "", err
	}

	token, err := GenerateActivationToken()
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(`
		INSERT INTO user_activation (id, user_id, token, expired_at)
		VALUES ($1, $2, $3, $4)
	`, uuid.New(), newID, token, time.Now().UTC().Add(24 * time.Hour))
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}
	err = helpers.SendActivationEmail(email, token)
	if err != nil {
		return "", err
	}
	return token, nil
}
func RegisterAdmin(db *sql.DB, email, password string) error {
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
		INSERT INTO users (id, email, password, role, is_active) 
		VALUES ($1, $2, $3, $4, $5)
	`, newID, email, string(hashedPassword), "Admin", true)
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

func RegisterDoctor(db *sql.DB, email, password string) error {
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
		INSERT INTO users (id, email, password, role, is_active) 
		VALUES ($1, $2, $3, $4, $5)
	`, newID, email, string(hashedPassword), "Doctor", true)
	if err != nil {
		return errors.New("failed to register doctor: " + err.Error())
	}

	_, err = tx.Exec(`
		INSERT INTO doctor (doctor_id, user_id) 
		VALUES ($1, $2)
	`, uuid.New(), newID)
	if err != nil {
		return errors.New("failed to create doctor: " + err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
func LoginUser(db *sql.DB, email, password string) (*structs.Users, string, error) {
	var user structs.Users

	err := db.QueryRow(`
		SELECT id, email, password, role, is_active
		FROM users WHERE email = $1
	`, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.IsActive)

	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, "", errors.New("account not activated")
	}

	token, err := GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
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

func GetAllCustomers(db *sql.DB, limit, offset int) (result []structs.Customer, err error) {
	sql := `
        SELECT  
			c.customer_id, c.status , u.id, u.email,
			u.role, u.is_active
		FROM customer c
		JOIN users u ON c.user_id = u.id
		LIMIT $1 OFFSET $2
    `

	rows, err := db.Query(sql, limit, offset)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var customer structs.Customer
		var user structs.Users

		err = rows.Scan(
			&customer.CustomerID, &customer.Status, &user.ID,
			&user.Email, &user.Role, &user.IsActive,
		)
		if err != nil {
			return
		}
		customer.User = &user
		result = append(result, customer)
	}

	return
}

func GetAllDoctors(db *sql.DB, limit, offset int) (result []structs.Doctor, err error) {
	sql := `
        SELECT  
			d.doctor_id, u.id, u.email,
			u.role, u.is_active
		FROM doctor d
		JOIN users u ON d.user_id = u.id
		LIMIT $1 OFFSET $2
    `

	rows, err := db.Query(sql, limit, offset)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var doctor structs.Doctor
		var user structs.Users

		err = rows.Scan(
			&doctor.DoctorID, &user.ID,
			&user.Email, &user.Role, &user.IsActive,
		)
		if err != nil {
			return
		}
		doctor.User = &user
		result = append(result, doctor)
	}

	return
}
func UpdateProfile(db *sql.DB, userID uuid.UUID, avatar []byte) error {
	var existing []byte
	err := db.QueryRow(`SELECT avatar FROM users WHERE id = $1`, userID).Scan(&existing)
	if err != nil {
		return errors.New("failed to fetch user: " + err.Error())
	}
	_, err = db.Exec(`
		UPDATE users SET avatar = $1 WHERE id = $2
	`, avatar, userID)
	if err != nil {
		return errors.New("failed to update profile: " + err.Error())
	}

	return nil
}
func GetUserByID(db *sql.DB, userID uuid.UUID) (structs.Users, error) {
    var user structs.Users
    err := db.QueryRow(`
        SELECT id, email, role, avatar, is_active 
        FROM users WHERE id = $1
    `, userID).Scan(&user.ID, &user.Email, &user.Role, &user.Avatar, &user.IsActive)
    if err != nil {
        return structs.Users{}, errors.New("failed to fetch user: " + err.Error())
    }
    return user, nil
}