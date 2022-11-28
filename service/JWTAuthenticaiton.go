package service

import (
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"math/rand"
	"sse_demo/util"
	"time"
)

// jwt service

type JWTService interface {
	GenerateToken(uuidStr string) string
	ValidateToken(token string) (*jwt.Token, error)
}
type AuthCustomClaims struct {
	UserID int32  `json:"user_id,omitempty"`
	UUID   string `json:"uuid,omitempty"`
	//K               string `json:"k"`
	Key            string `json:"key"`
	Iv             string `json:"iv"`
	Quality        string `json:"quality"`
	IsLcVideoAdmin bool   `json:"is_lc_video_admin"`
	IsSup          bool   `json:"is_sup"`
	jwt.StandardClaims
}

type jwtServices struct {
	secretKey string
}

// auth-jwt

func JWTAuthService(secretKey string) JWTService {
	return &jwtServices{
		secretKey: getSecretKey(secretKey),
	}
}

func getSecretKey(secretKey string) string {
	secret := viper.GetString(secretKey)
	if secret == "" {
		secret = "lc-bc-zp^95mf083s1rfoe15v9wouxv42hb04e1i4***YJz"
	}
	return secret
}

func getK() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	hexK := hex.EncodeToString(bytes)
	return hexK
}

func ParseK(k string) ([]byte, []byte) {
	util.MyLogger.Info(k)
	keyStr, ivStr := k[:32], k[32:]
	keyHex, _ := hex.DecodeString(keyStr)
	ivHex, _ := hex.DecodeString(ivStr)
	return keyHex, ivHex
}

func getKeyAndIv() (string, string) {
	keyBytes, ivBytes := make([]byte, 16), make([]byte, 16)
	rand.Read(keyBytes)
	rand.Read(ivBytes)
	hexKey, hexIv := hex.EncodeToString(keyBytes), hex.EncodeToString(ivBytes)
	return hexKey, hexIv
}

func ParseKeyAndIv(key, iv string) ([]byte, []byte) {
	keyHex, _ := hex.DecodeString(key)
	ivHex, _ := hex.DecodeString(iv)
	return keyHex, ivHex
}

func (service *jwtServices) GenerateToken(uuidStr string) string {

	key, iv := getKeyAndIv()
	claims := &AuthCustomClaims{
		1,
		uuidStr,
		key,
		iv,
		"",
		false,
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// encoded string
	t, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (service *jwtServices) ValidateToken(encodedToken string) (*jwt.Token, error) {
	var tempClaims AuthCustomClaims
	return jwt.ParseWithClaims(encodedToken, &tempClaims, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("Invalid token:", token.Header["alg"])

		}
		return []byte(service.secretKey), nil
	})

}
