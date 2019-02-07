package keystore

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	crypto "github.com/libp2p/go-libp2p-crypto"
)

// TokenClaims represents a token
type TokenClaims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Exp   int64  `json:"exp"`
	Iat   int64  `json:"iat"`
}

func NewTokenClaims(name string, email string) TokenClaims {
	return TokenClaims{
		Name:  name,
		Email: email,
		Exp:   time.Now().Add(time.Hour).Unix(),
		Iat:   time.Now().Unix(),
	}
}

func (c TokenClaims) Valid() error {
	return nil
}

// A JWT Token.  Different fields will be used depending on whether you're
// creating or parsing/verifying a token.
type Token struct {
	Raw         string                 // The raw token.  Populated when you Parse a token
	Header      map[string]interface{} // The first segment of the token
	TokenClaims TokenClaims            // The second segment of the token
	Signature   string                 // The third segment of the token.  Populated when you Parse a token
	Valid       bool                   // Is the token valid?  Populated when you Parse/Verify a token
}

func NewWithClaims(claims TokenClaims) *Token {
	return &Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": "HS256",
		},
		TokenClaims: claims,
	}
}

func (t *Token) GetToken(key crypto.PrivKey) (string, error) {
	var sstr string
	var sig []byte
	var err error
	if sstr, err = t.CreateTokenParts(); err != nil {
		return "", err
	}
	fmt.Println("sstr: ")
	fmt.Println(sstr)
	fmt.Println(key)
	if sig, err = key.Sign([]byte(sstr)); err != nil {
		return "", err
	}
	return strings.Join([]string{sstr, hex.EncodeToString(sig)}, "."), nil
}

func (t *Token) CreateTokenParts() (string, error) {
	var err error
	parts := make([]string, 2)
	for i, _ := range parts {
		var jsonValue []byte
		if i == 0 {
			if jsonValue, err = json.Marshal(t.Header); err != nil {
				return "", err
			}
		} else {
			if jsonValue, err = json.Marshal(t.TokenClaims); err != nil {
				return "", err
			}
		}

		parts[i] = EncodeSegment(jsonValue)
	}
	return strings.Join(parts, "."), nil
}

// Encode JWT specific base64url encoding with padding stripped
func EncodeSegment(seg []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(seg), "=")
}

func VerifyToken(pubKey crypto.PubKey, token string) (bool, error) {
	tokenParts := strings.Split(token, ".")[:2]
	signature := strings.Split(token, ".")[2]
	data := strings.Join(tokenParts, ".")
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false, err
	}
	return pubKey.Verify([]byte(data), sigBytes)
}

// func CreateJSONToken(tokenData string) ([]byte, error) {
// 	token := TokenJSON{
// 		TokenData: tokenData,
// 	}
// 	data, err := json.MarshalIndent(&token, "", "  ")
// 	if err != nil {
// 		return data, err
// 	}
// 	return data, nil
// }

// func DecryptJSONToken(tokenData string) (string, error) {
// 	encjson := TokenJSON{}
// 	err := json.Unmarshal([]byte(tokenData), &encjson)
// 	if err != nil {
// 		return "", err
// 	}
// 	return encjson.TokenData, nil
// }

func CreateToken(privKey crypto.PrivKey) (string, error) {
	token := NewWithClaims(NewTokenClaims("Test", "aa@bb.com"))
	tokenString, err := token.GetToken(privKey)
	if err != nil {
		fmt.Println("There was an error: ", err)
		return "", err
	}
	fmt.Println("tokenString:")
	fmt.Println(tokenString)
	return tokenString, nil
}
