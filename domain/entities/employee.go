package entities

import (
	"encoding/hex"
	"time"

	"github.com/bekha-io/openbank/domain/types"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Employee struct {
	ID         types.EmployeeID `json:"id"`
	Email      string           `json:"email"`
	Password   string           `json:"-"`
	FirstName  string           `json:"first_name"`
	LastName   string           `json:"last_name"`
	MiddleName string           `json:"middle_name"`
	ImageUrl   string           `json:"image_url"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

// NewEmployee creates a new Employee based on the given parameters. Password is hashed as result
func NewEmployee(email, password, firstName, lastName, middleName string) *Employee {
	return &Employee{
		ID:         types.EmployeeID(uuid.New()),
		Email:      email,
		Password:   hashPassword(password),
		FirstName:  firstName,
		LastName:   lastName,
		MiddleName: middleName,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}

func hashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hex.EncodeToString(hashedPassword)
}

func (e *Employee) IsPasswordCorrect(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(e.Password), []byte(password)) == nil
}
