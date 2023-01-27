package helpers

import (
	"docker_test/src/config"
	"docker_test/src/db"
	"encoding/json"
	"fmt"

	"docker_test/structs"
	DB "docker_test/structs/db"
	"net/http"
	"time"

	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func validateTokenHS256(token string) (*structs.CustomClaims, error) {

	// line verify
	url := fmt.Sprintf("https://api.line.me/v2/profile?access_token=%s", token)
	req, _ := http.NewRequest("GET", url, nil)
	// add header bearer
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	if res, _ := http.DefaultClient.Do(req); res.StatusCode == http.StatusOK {
		// read body
		var lineProfile structs.LineProfile
		if err := json.NewDecoder(res.Body).Decode(&lineProfile); err == nil {
			claims := &structs.CustomClaims{
				UserID: lineProfile.UserID,
				Email:  lineProfile.Email,
				StandardClaims: jwt.StandardClaims{
					Audience:  "docker_test",
					Issuer:    lineProfile.DisplayName,
					ExpiresAt: time.Now().Unix() + config.TOKENEXPIRETIME,
				},
			}
			// return claims
			return claims, nil

		}

	}

	// custom verify
	claims := &structs.CustomClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		if err.Error() == "signature is invalid" && claims.Issuer != "" {
			claims.UserID = claims.Issuer
			return claims, nil
		}
		return nil, err
	}
	return claims, nil
}

func GenTokenHS256(UserID, Username string) (string, *structs.CustomClaims, error) {
	claims := &structs.CustomClaims{
		UserID:   UserID,
		Username: Username,
		StandardClaims: jwt.StandardClaims{
			Audience:  "docker_test",
			Issuer:    Username,
			ExpiresAt: time.Now().Unix() + config.TOKENEXPIRETIME,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", nil, err
	}
	// is user exist
	if err := db.Con.Where("uid = ?", claims.UserID).First(&DB.Users{}).Error; err != nil {
		return "", nil, err
	}

	return tokenString, claims, nil
}

func DecodeTokenHS256(token string) (*structs.CustomClaims, error) {
	// remove barer
	token = strings.Replace(token, "Bearer ", "", 1)
	token = strings.TrimSpace(token)
	// decode
	claims := &structs.CustomClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func GetTokenFormHeader(c echo.Context) string {
	// Get token from the request
	auth := c.Request().Header.Get("Authorization")
	token := strings.Replace(auth, "Bearer ", "", 1)
	token = strings.TrimSpace(token)
	return token
}

func ResetTokenByContext(c echo.Context) error {
	token := GetTokenFormHeader(c)
	claims, err := validateTokenHS256(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Unauthorized",
		})
	}
	// gen new token
	newToken, _, err := GenTokenHS256(claims.UserID, claims.Username)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Unauthorized",
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": newToken,
	})
}

func ResetTokenByUserIdUsername(ID *uint, UserID, Username string, c echo.Context) (echo.Context, error) {
	// gen new token
	newToken, _, err := GenTokenHS256(UserID, Username)
	if err != nil {
		return c, err
	}
	claims, err := validateTokenHS256(newToken)
	if err != nil {
		return c, err
	}

	c.Response().Header().Set("Authorization", "Bearer "+newToken)
	c.Set("id", ID)
	c.Set("user", claims)
	return c, nil
}
