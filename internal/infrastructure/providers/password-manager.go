// Package providers provides components that focus on providing common services.
package providers

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

// A PasswordManager allows generating and verifying passwords.
// PasswordManager uses the bcrypt algorithm internally to hash and
// compare passwords.
type PasswordManager struct {
	// Cost specifies the computational cost for hashing passwords
	// using the bcrypt algorithm.
	cost int
}

// NewPasswordManager creates a new PasswordManager.
func NewPasswordManager(cost int) *PasswordManager {
	return &PasswordManager{
		cost: cost,
	}
}

// validateRestrictions validates that the password meets the requirements
// to be a valid password:
// - Length <= 72 bytes
// - Cannot contain null characters: \x00
func (s *PasswordManager) validateRestrictions(password string) error {
	passwordBytes := []byte(password)
	if len(passwordBytes) > 72 {
		return fmt.Errorf(
			"invalid password: password is %d bytes, exceeds 72 bytes",
			len(passwordBytes),
		)
	}
	if strings.ContainsRune(password, '\x00') {
		return fmt.Errorf("invalid password: password contains null characters")
	}
	return nil
}

// HashPassword returns the hashed password. The password length must be less than or
// equal to 72 bytes and must not contain null characters; otherwise, a non-nil error will
// be returned.
func (s *PasswordManager) HashPassword(password string) (string, error) {
	if err := s.validateRestrictions(password); err != nil {
		return "", err
	}
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", fmt.Errorf("password hashing error: %w", err)
	}
	return string(hashedPasswordBytes), nil
}

// CheckPassword compares whether a hashed password is a valid hash for a plaintext
// password.
func (s *PasswordManager) CheckPassword(hashedPassword string, password string) bool {
	if err := s.validateRestrictions(password); err != nil {
		return false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}
	return true
}