package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/harshgupta9473/sevak_backend/models"
	"github.com/harshgupta9473/sevak_backend/utils"
	"github.com/joho/godotenv"
)

type AuthMiddleWare func(http.ResponseWriter, *http.Request) http.HandlerFunc

func LoadKeys() (string, string, string, error) {
	 err := godotenv.Load()
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}
	jwtLogin := os.Getenv("secretKeyForLoginJWT")
	jwtRequest := os.Getenv("secretKeyForRequests")
	jwtRefresh := os.Getenv("secretKeyForRefreshJWT")

	return jwtLogin, jwtRequest, jwtRefresh, nil
}

func CreateLoginJWT(user *models.User, role string) (string, error) {
	secretkey, _, _, err := LoadKeys()
	if err != nil {
		
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authid":     user.Mobile,
		"user":       role,
		"exp":        time.Now().Add(15 * time.Minute).Unix(),
		"created_at": user.CreatedAt,
	})
	tokenString, err := token.SignedString([]byte(secretkey))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return tokenString, nil
}

func CreateRefreshJWT(user *models.User, role string) (string, error) {
	_, _, secretkey, err := LoadKeys()
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authid":     user.Mobile,
		"user":       role,
		"exp":        time.Now().Add(7 * 24 * time.Hour).Unix(),
		"created_at": user.CreatedAt,
	})
	tokenString, err := token.SignedString([]byte(secretkey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string, login bool) (*jwt.Token, error) {
	loginkey, _, refreshkey, err := LoadKeys()

	if err != nil {
		return nil, err
	}
	var secretkey string
	if login {
		secretkey = loginkey
	} else {
		secretkey = refreshkey
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretkey), nil
	})
}




func AuthMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
		const userContextKey contextKey = "userInfo"
		// Extract the access token from cookies
		accessCookie, err := r.Cookie("token") // Assuming the cookie is named "token"
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Authorization cookie is missing", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Error retrieving cookie: "+err.Error(), http.StatusUnauthorized)
			return
		}

		accessTokenString := accessCookie.Value

		// Validate the access token
		accessToken, err := ValidateToken(accessTokenString, true)
		if err != nil {
			http.Error(w, "Invalid access token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Check if the access token is valid
		if !accessToken.Valid {
			// Token is not valid, check for expiration
			if claims, ok := accessToken.Claims.(jwt.MapClaims); ok {
				exp := int64(claims["exp"].(float64))
				role:=claims["role"].(string)
				if role=="admin"{
					http.Error(w,"not allowed",http.StatusForbidden)
					return
				}
				if time.Now().Unix() > exp {
					// Attempt to use the refresh token to issue a new access token
					refreshCookie, err := r.Cookie("refresh_token") 
					if err != nil {
						if err == http.ErrNoCookie {
							http.Error(w, "Refresh token cookie is missing", http.StatusUnauthorized)
							return
						}
						http.Error(w, "Error retrieving refresh cookie: "+err.Error(), http.StatusUnauthorized)
						return
					}

					refreshTokenString := refreshCookie.Value

					// Validate the refresh token
					refreshToken, err := ValidateToken(refreshTokenString, false) // Use the refresh token validation logic
					if err != nil || !refreshToken.Valid {
						http.Error(w, "Invalid refresh token: "+err.Error(), http.StatusUnauthorized)
						return
					}

					// Check if the refresh token is valid and not expired
					if refreshClaims, ok := refreshToken.Claims.(jwt.MapClaims); ok {
						role:=claims["role"].(string)
				if role=="admin"{
					http.Error(w,"not allowed",http.StatusForbidden)
					return
				}
						refreshExp := int64(refreshClaims["exp"].(float64))
						if time.Now().Unix() > refreshExp { // Refresh token is expired
							http.Error(w, "Refresh token is expired", http.StatusUnauthorized)
							return
						}

						// Extract user information from the refresh token claims
						userMobile := refreshClaims["authid"].(string)
						userRole := refreshClaims["user"].(string)

						user, err := utils.AccessUSERSByMobileNumber(userMobile)
						if err != nil {
							http.Error(w, "unauthorised", http.StatusUnauthorized)
							return
						}
						if !utils.IsRoleAlloted(user.Roles, userRole) {
							http.Error(w, "unauthorised", http.StatusUnauthorized)
							return
						}
						// Create a new access token
						newAccessToken, err := CreateLoginJWT(user, userRole)
						if err != nil {
							http.Error(w, "Error creating new access token: "+err.Error(), http.StatusInternalServerError)
							return
						}

						// Set the new access token in a cookie
						utils.SetAuthTokenCookie(w, newAccessToken, "authtoken")
						fmt.Fprintln(w, "New access token generated.")
					} else {
						http.Error(w, "Invalid refresh token claims", http.StatusUnauthorized)
						return
					}
				}
			}
			http.Error(w, "Invalid access token", http.StatusUnauthorized)
			return
		}
		claims,_:= accessToken.Claims.(jwt.MapClaims);
		role:=claims["role"].(string)
				if role=="admin"{
					http.Error(w,"not allowed",http.StatusForbidden)
					return
				}
		userInfo:=map[string]interface{}{
			"number":claims["authid"].(string),
			"role":claims["user"].(string),
			"created_at":claims["created_at"].(time.Time),
		}
		ctx:=context.WithValue(r.Context(),userContextKey,userInfo)
		next.ServeHTTP(w,r.WithContext(ctx))
		
	})
}
