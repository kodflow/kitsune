package handler

import (
	"io"
	"net/http"

	"github.com/kodmain/kitsune/src/internal/core/server/router"
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
)

func HTTPHandler(w http.ResponseWriter, r *http.Request) {
	req, res := transport.New()
	req.Method = r.Method

	if r.Method == "POST" || r.Method == "PATCH" || r.Method == "PUT" {
		// Lire le corps de la requête
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, "Erreur lors de la lecture de la requête", http.StatusBadRequest)
			return
		}

		req.Body = body
	}

	// Traiter la requête
	if err := router.Resolve(req, res); err != nil {
		http.Error(w, "Erreur de traitement", http.StatusInternalServerError)
		return
	}

	// Écrire la réponse
	w.Header().Set("request-id", req.Id)
	w.WriteHeader(int(res.Status))
	w.Write(res.Body)
}
