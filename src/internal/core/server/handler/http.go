package handler

import (
	"io"
	"net/http"

	"github.com/kodmain/kitsune/src/internal/core/server/router"
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
)

// HTTPHandler handles HTTP requests and sends back HTTP responses.
// It processes incoming HTTP requests, creates a corresponding transport request,
// and uses the router to generate a response. The handler deals with various HTTP methods,
// reads request bodies if necessary, and writes back responses including headers and status codes.
//
// Parameters:
// - w: http.ResponseWriter Response writer to send back the HTTP response.
// - r: *http.Request The incoming HTTP request to be processed.
func HTTPHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize a new transport request and response
	req, res := transport.New()
	req.Method = r.Method

	// Read the request body for specific HTTP methods
	if r.Method == "POST" || r.Method == "PATCH" || r.Method == "PUT" {
		// Read the request body
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			// Handle errors in reading the request body
			http.Error(w, "Erreur lors de la lecture de la requête", http.StatusBadRequest)
			return
		}

		req.Body = body
	}

	// Process the request using the router
	if err := router.Resolve(req, res); err != nil {
		// Handle errors in processing the request
		http.Error(w, "Erreur de traitement", http.StatusInternalServerError)
		return
	}

	// Write the response back to the client
	w.Header().Set("request-id", req.Id)
	w.WriteHeader(int(res.Status))
	w.Write(res.Body)
}
