package requests

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Request struct {
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate validates a parsed request using the go-playground validator
func (m *Request) Validate(req interface{}) error {
	// Validate the request using the initialized validator
	if err := validate.Struct(req); err != nil {
		// If validation fails, construct a human-readable error message
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			fieldName := e.Field()
			tagName := e.Tag()
			errorMessage := e.Param()

			// Customize error messages based on validation tags
			switch tagName {
			case "required":
				validationErrors = append(validationErrors, fmt.Sprintf("%s is required", fieldName))
			case "min":
				validationErrors = append(validationErrors, fmt.Sprintf("%s must be at least %s characters", fieldName, errorMessage))
			case "max":
				validationErrors = append(validationErrors, fmt.Sprintf("%s cannot be longer than %s characters", fieldName, errorMessage))
			default:
				validationErrors = append(validationErrors, fmt.Sprintf("%s is not valid", fieldName))
			}
		}

		// Join validation error messages
		errorMessage := strings.Join(validationErrors, ", ")

		// Return an error with the formatted validation errors
		return fiber.NewError(fiber.StatusBadRequest, errorMessage)
	}

	// If validation succeeds, return nil
	return nil
}

// ParseAndValidate parses and validates a request
func (m *Request) ParseAndValidate(c *fiber.Ctx, req interface{}) error {
	// Check content type and parse accordingly
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Validate the request
	if err := m.Validate(req); err != nil {
		return err
	}

	// Return nil if parsing and validation succeed
	return nil
}
