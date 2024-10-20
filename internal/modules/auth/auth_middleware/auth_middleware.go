package auth_middleware

import (
	"context"
	"github.com/golang-jwt/jwt"
	"gophKeeper/internal/logger"
	"gophKeeper/internal/modules/auth/auth_services/auth_jwt_service"
	"net/http"
	"strings"
)

func Authentication(next http.Handler, secretKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtBearer := r.Header.Get("Authorization")
		logger.Log.Infoln(jwtBearer)
		strArr := strings.Split(jwtBearer, " ")
		if len(strArr) != 2 {
			logger.Log.Errorln("jwt bearer format error")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		logger.Log.Infoln(strArr[1])
		token, err := auth_jwt_service.ParseAndValidateToken(strArr[1], secretKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			logger.Log.Infoln(token)
		}

		ctx := context.Background()
		ctx = context.WithValue(ctx, "token", token)
		r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
