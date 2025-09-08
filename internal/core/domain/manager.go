package domain

import (
	"golang.org/x/crypto/bcrypt"
)

type Manager struct {
	Id       string
	Email    string
	Password string
}

func (m *Manager) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	m.Password = string(hashedPassword)
	return nil
}

func (m *Manager) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(password))
}
