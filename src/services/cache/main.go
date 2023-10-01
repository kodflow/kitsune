package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Créer une nouvelle instance de l'application Fiber
	app := fiber.New()

	// Route de base pour vérifier que le serveur fonctionne
	app.Get("/", func(c *fiber.Ctx) error {
		// Input: *fiber.Ctx qui représente le contexte de la requête/réponse
		// Output: error, une éventuelle erreur à renvoyer
		// Objectif: Renvoyer un message simple pour confirmer que le serveur est opérationnel
		return c.SendString("Hello, World!")
	})

	// Démarrer le serveur sur le port 3000
	log.Fatal(app.Listen(":9998"))
}
