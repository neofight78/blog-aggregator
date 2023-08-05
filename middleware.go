package main

import (
	"com.github/neofight78/blog-aggregator/internal/database"
	"net/http"
	"strings"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		fields := strings.Fields(authorization)

		if len(fields) != 2 || fields[0] != "ApiKey" {
			respondWithError(w, http.StatusUnauthorized, "Missing or invalid authorization header")
			return
		}

		user, err := cfg.Queries.GetUser(r.Context(), fields[1])
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		handler(w, r, user)
	}
}
