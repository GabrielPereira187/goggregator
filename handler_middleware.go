package main

import (
	"net/http"

	"github.com/GabrielPereira187/blog-aggregator/internal/auth"
)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        token, err := auth.GetApiKey(r.Header)
        if err != nil {
            respondWithError(w, 403, "Unauthorized")
            return
        }

        handler(w, r, token)
    }
}