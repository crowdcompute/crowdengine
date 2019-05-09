// Copyright 2018 The crowdcompute:crowdengine Authors
// This file is part of the crowdcompute:crowdengine library.
//
// The crowdcompute:crowdengine library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The crowdcompute:crowdengine library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the crowdcompute:crowdengine library. If not, see <http://www.gnu.org/licenses/>.

package keystore

import (
	"fmt"
	"log"
	"time"

	"github.com/crowdcompute/crowdengine/crypto"
	jwt "github.com/dgrijalva/jwt-go"
)

// TokenClaims represents a token
type TokenClaims struct {
	// TODO: See what do we actually need as claims
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	jwt.StandardClaims
}

// NewTokenClaims creates and returns new Token Claims
func NewTokenClaims(name string, email string) *TokenClaims {
	sc := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
		IssuedAt:  time.Now().Unix(),
	}
	return &TokenClaims{
		Name:           name,
		Email:          email,
		StandardClaims: sc,
	}
}

// Valid determines if the token is invalid for any supported reason
// This is being checked when we are Parsing the token
func (c *TokenClaims) Valid() error {
	return c.StandardClaims.Valid()
}

// NewToken creates a new JWT token and fills the Raw of the jwt.Token
func NewToken(key []byte, tClaims *TokenClaims) (*jwt.Token, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, tClaims)
	tokenString, err := tok.SignedString(key)
	if err != nil {
		log.Println("There was an error getting signed token: ", err)
		return nil, err
	}
	tok.Raw = tokenString
	return tok, nil
}

func HashToken(rawToken string) string {
	b := crypto.Sha256Hash([]byte(rawToken)).Sum(nil)
	return string(b)
}

// VerifyToken checks if the token's data is valid
func VerifyToken(rawToken string, key []byte) (bool, error) {
	parser := new(jwt.Parser)
	// Checks if the claims are valid as well
	token, err := parser.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return key, nil
	})
	if err != nil {
		return false, err
	}
	return token.Valid, err
}

// VerifyTokenWithClaims checks if the token's data is valid
func VerifyTokenWithClaims(t *jwt.Token, key []byte) (bool, error) {
	parser := new(jwt.Parser)
	// Checks if the claims are valid as well
	token, err := parser.ParseWithClaims(t.Raw, t.Claims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return key, nil
	})
	return token.Valid, err
}
