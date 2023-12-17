package promise

import (
	"testing"

	"github.com/kodflow/kitsune/src/internal/core/server/transport/proto/generated"
	"github.com/stretchr/testify/assert"
)

// TestPromiseLifecycle teste la création, l'ajout et la résolution d'une promesse.
func TestPromiseLifecycle(t *testing.T) {
	// Création d'une promesse
	callback := func(responses ...*generated.Response) {
		// Logique de rappel
	}
	p, err := Create(callback)
	assert.NoError(t, err, "Erreur lors de la création de la promesse")

	// Ajout d'une réponse
	req := &generated.Request{}
	p.Add(req)

	// Test de la résolution
	res := &generated.Response{}
	p.Resolve(res)
	assert.Equal(t, 1, len(p.responses), "Le nombre de réponses n'est pas correct")
	assert.True(t, p.Closed, "La promesse n'est pas fermée après la résolution")
}

// TestRepositoryLifecycle teste la création et la recherche de promesses dans le dépôt.
func TestRepositoryLifecycle(t *testing.T) {
	p, err := Create(func(responses ...*generated.Response) {
		// do it
	})
	assert.NoError(t, err, "Erreur lors de la création de la promesse")

	foundPromise, err := Find(p.Id)
	assert.NoError(t, err, "Erreur lors de la recherche de la promesse")
	assert.Equal(t, p.Id, foundPromise.Id, "Les IDs des promesses ne correspondent pas")
}
