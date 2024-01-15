package http_test

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/kodflow/kitsune/src/internal/core/server/protocols/http"
	"github.com/kodflow/kitsune/src/internal/core/server/transport/proto/generated"
	"github.com/stretchr/testify/assert"
)

func generateTwoDistinctRandomNumbers() (string, string) {
	min, max := 2000, 65000

	firstNumber := rand.Intn(max-min+1) + min
	secondNumber := rand.Intn(max-min+1) + min

	// Régénérer le second nombre jusqu'à ce qu'il soit différent du premier.
	for firstNumber == secondNumber {
		secondNumber = rand.Intn(max-min+1) + min
	}

	return strconv.Itoa(firstNumber), strconv.Itoa(secondNumber)
}

func TestHTTPClient(t *testing.T) {
	p1, p2 := generateTwoDistinctRandomNumbers()
	// Création d'un serveur HTTP de test
	server := setupHTTPServer(p1, p2)
	server.Start()
	defer server.Stop()

	// Initialisation du client HTTP
	client := http.NewHTTPClient()

	req := &generated.Request{
		Method:   "GET",
		Endpoint: "http://127.0.0.1:" + p1,
	}

	t.Run("GET Request", func(t *testing.T) {
		res := client.Send(req)
		assert.Equal(t, uint32(404), res.Status)
	})

	req = &generated.Request{
		Method:   "GET",
		Endpoint: "https://127.0.0.1:" + p2,
	}

	t.Run("GET Request", func(t *testing.T) {
		res := client.Send(req)
		assert.Equal(t, uint32(404), res.Status)
	})
}
