package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/cristalhq/jwt/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var (
	ErrEntityNotFound    = NewServiceError(404, "Entity Not Found")
	ErrInvalidIdentifier = NewServiceError(400, "Invalid Identifier")
)

type DataEnvelope struct {
	Data interface{} `json:"data"`
}

type ErrorEnvelope struct {
	Error interface{} `json:"error"`
}

type ServiceError struct {
	Title   string `json:"title"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewServiceError(status int, message string) *ServiceError {
	return &ServiceError{
		Title:   http.StatusText(status),
		Status:  status,
		Message: message,
	}
}

// TODO: Allow usage of name or UUID for singular endpoints.

var isAlphaNumeric = regexp.MustCompile(`^[0-9a-zA-Z]+$`).MatchString

func API(db *sqlx.DB) *fiber.App {
	api := fiber.New()

	// Create a new Signer and Builder to create JWTs using the HMAC algorithm.
	jwtKeyVariable := "JWT_KEY"
	jwtKey := os.Getenv(jwtKeyVariable)
	if jwtKey == "" {
		log.Fatal().Msgf("Missing environment variable: %s", jwtKeyVariable)
	}

	signer, err := jwt.NewSignerHS(jwt.HS512, []byte(jwtKey))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create JWT signer")
	}

	builder := jwt.NewBuilder(signer)

	api.Get("/queues", func(c *fiber.Ctx) error {
		queues := make([]Queue, 0)
		if err := db.Select(&queues, `SELECT * FROM queues`); err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(queues)
	})

	api.Post("/queues", func(c *fiber.Ctx) error {
		// Parse request body.
		entity := new(Queue)
		if err := c.BodyParser(entity); err != nil {
			return err
		}

		// TODO: Validate request body.

		// Assign new UUID.
		entity.ID = uuid.NewString()

		// Insert new entity.
		_, err := db.NamedExec(`
            INSERT
            INTO queues (id, name, owner, title, description, number)
            VALUES (:id, :name, :owner, :title, :description, :number)
        `, entity)
		if err != nil {
			// TODO: Handle unique constraint errors gracefully.
			return err
		}

		// Create new JWT token for authentication.
		token, err := builder.Build(&jwt.RegisteredClaims{
			ID:      fmt.Sprintf("/api/queues/%s", entity.ID),
			Subject: entity.Owner,
		})
		if err != nil {
			return err
		}
		entity.Token = token.String()

		return c.Status(fiber.StatusCreated).JSON(DataEnvelope{
			Data: entity,
		})
	})

	api.Get("/queues/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")

		// Prevent SQL injection.
		if !isAlphaNumeric(name) {
			return c.Status(ErrInvalidIdentifier.Status).JSON(ErrorEnvelope{
				Error: *ErrInvalidIdentifier,
			})
		}

		entity := new(Queue)
		if err := db.Get(entity, "SELECT * FROM queues WHERE name = $1", name); err != nil {
			return c.Status(ErrEntityNotFound.Status).JSON(ErrorEnvelope{
				Error: *ErrEntityNotFound,
			})
		}

		return c.Status(fiber.StatusOK).JSON(entity)
	})

	api.Put("/queues/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")

		// Prevent SQL injection.
		if !isAlphaNumeric(name) {
			return c.Status(ErrInvalidIdentifier.Status).JSON(ErrorEnvelope{
				Error: *ErrInvalidIdentifier,
			})
		}

		// Parse request body.
		entity := new(Queue)
		if err := c.BodyParser(entity); err != nil {
			return err
		}
		entity.Name = name

		_, err := db.NamedExec(`
            UPDATE queues
            SET title = :title,
                description = :description,
                number = :number
            WHERE name = :name
        `, entity)
		if err != nil {
			return err
		}

		// TODO: Have someone review this. This doesn't really seem
		// like the right way to check if the query succeeded.
		if err := db.Get(entity, "SELECT * FROM queues WHERE name = $1", name); err != nil {
			return c.Status(ErrEntityNotFound.Status).JSON(ErrorEnvelope{
				Error: *ErrEntityNotFound,
			})
		}

		return c.Status(fiber.StatusOK).JSON(entity)
	})

	api.Delete("/queues/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")

		// Prevent SQL injection.
		if !isAlphaNumeric(name) {
			return c.Status(ErrInvalidIdentifier.Status).JSON(ErrorEnvelope{
				Error: *ErrInvalidIdentifier,
			})
		}

		res, err := db.Exec("DELETE FROM queues WHERE name = $1", name)
		if err != nil {
			return err
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return c.Status(ErrEntityNotFound.Status).JSON(ErrorEnvelope{
				Error: *ErrEntityNotFound,
			})
		}

		if rowsAffected != 1 {
			return fmt.Errorf("expected 1 row to be affected, got %d", rowsAffected)
		}

		return c.SendStatus(fiber.StatusNoContent)
	})

	return api
}
