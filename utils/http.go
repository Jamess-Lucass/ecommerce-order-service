package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Jamess-Lucass/ecommerce-order-service/middleware"
	"github.com/golang-jwt/jwt/v4"
)

func HttpGet(uri string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	claims := middleware.Claim{RegisteredClaims: jwt.RegisteredClaims{Subject: "order-service", ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour))}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	secret := []byte(os.Getenv("JWT_SECRET"))
	jwt, err := token.SignedString(secret)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(string(body))
	}

	return res.Body, nil
}
