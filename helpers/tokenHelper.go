package helpers

import (
	"gastrono-go/database"
	"github.com/dgrijalva/jwt-go"
	"os"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	jwt.StandardClaims
}

var userCollection = database.OpenCollection(database.Client, "user")

var SECRET_KEY = os.Getenv("SECRET_KEY")

func GenerateAllTokens() {}

func UpdateAllTokens() {}

func ValidateToken() {}
