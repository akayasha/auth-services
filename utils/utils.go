package utils

import (
	"auth-services/models"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

// Hashing Password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Check Hashing(password)
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Generate JWT
func GenerateJWT(uuid string, nip string, role models.Role) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = uuid
	claims["role"] = string(role)
	claims["nip"] = nip
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}

// GenerateOTP generates a random OTP
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(999999-100000) + 100000
	return fmt.Sprintf("%d", otp)
}

func ValidateStruct(input interface{}, requiredFields ...string) string {
	missingFields := []string{}
	val := reflect.ValueOf(input)

	for _, field := range requiredFields {
		f := val.FieldByName(field)
		if !f.IsValid() || (f.Kind() == reflect.String && strings.TrimSpace(f.String()) == "") {
			missingFields = append(missingFields, field)
		}
	}

	if len(missingFields) > 0 {
		return fmt.Sprintf("Missing or empty fields: %v", strings.Join(missingFields, ", "))
	}
	return ""
}
