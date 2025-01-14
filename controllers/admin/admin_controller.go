package adminpanel

import (
	"context"
	mongoSetup "em_backend/configs/mongo"
	"encoding/json"
	"fmt"
	"strings"

	"time"

	commonutils "em_backend/library/common"
	dbModel "em_backend/models/db"
	common_responses "em_backend/responses/common"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/google/uuid"
)

func CreateEvent(ctx *fiber.Ctx) error {
	// Parse request body
	var requestData dbModel.Event
	fmt.Print("vate")
	if err := ctx.BodyParser(&requestData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing request data",
			"status":  "400 Bad Request",
		})
	}

	// Validate required fields
	if strings.TrimSpace(requestData.EventName) == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Event name is required",
			"status":  "400 Bad Request",
		})
	}

	// if requestData.RegistrationLimit <= 0 {
	// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"message": "Registration limit must be greater than 0",
	// 		"status":  "400 Bad Request",
	// 	})
	// }

	if requestData.EventDate < time.Now().Unix() {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Event date cannot be in the past",
			"status":  "400 Bad Request",
		})
	}

	// if requestData.PaymentType == "paid" && len(requestData.RegistrationData) == 0 {
	// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"message": "Registration data must be provided for paid events",
	// 		"status":  "400 Bad Request",
	// 	})
	// }

	// Generate unique IDs
	eventID, err := uuid.NewRandom()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error generating event ID",
			"status":  "500 Internal Server Error",
		})
	}
	requestData.UniqueId = eventID.String()
	requestData.CreatedAt = time.Now().Unix()
	requestData.Status = "active"

	formID, err := uuid.NewRandom()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error generating registration form ID",
			"status":  "500 Internal Server Error",
		})
	}
	requestData.RegistrationDetailsFormId = formID.String()

	// Connect to MongoDB
	db, eventCollection, err := mongoSetup.ConnectMongo("events")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error connecting to MongoDB",
			"status":  "500 Internal Server Error",
		})
	}
	defer db.Client().Disconnect(context.Background())

	registerFormCollection := db.Collection("registerForm")

	// Insert event data (excluding RegistrationData and RegistrationForm)
	eventData := requestData
	eventData.RegistrationData = nil
	eventData.RegistrationForm = nil
	_, err = eventCollection.InsertOne(ctx.Context(), eventData)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to insert event",
			"status":  "500 Internal Server Error",
		})
	}

	// Insert registration form data
	if len(requestData.RegistrationForm) > 0 {
		registrationForm := dbModel.RegistrationForm{
			RegistrationFormId:     formID.String(),
			EventId:                eventID.String(),
			RegistrationFormFields: requestData.RegistrationForm,
		}
		_, err = registerFormCollection.InsertOne(ctx.Context(), registrationForm)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to insert registration form",
				"status":  "500 Internal Server Error",
			})
		}
	}

	// Return success response
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Event created successfully",
		"status":  "201 Created",
		"data": fiber.Map{
			"eventId":            requestData.UniqueId,
			"registrationFormId": requestData.RegistrationDetailsFormId,
		},
	})
}

func EditEvent(ctx *fiber.Ctx) error {
	// Parse request body
	var requestData dbModel.Event
	err := json.Unmarshal(ctx.Body(), &requestData)
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error parsing request data",
			Status:  "400 Bad Request",
		}))
	}

	// Validate the event ID
	eventID := requestData.UniqueId
	if eventID == "" {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Event ID is required",
			Status:  "400 Bad Request",
		}))
	}

	// Connect to MongoDB
	db, col, err := mongoSetup.ConnectMongo("events")
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error connecting to MongoDB",
			Status:  "500 Internal Server Error",
		}))
	}
	defer db.Client().Disconnect(context.TODO())

	// Find the event in the database
	var existingEvent dbModel.Event
	err = col.FindOne(ctx.Context(), bson.M{"uniqueId": eventID}).Decode(&existingEvent)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
				Message: "Event not found",
				Status:  "404 Not Found",
			}))
		}
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error fetching event",
			Status:  "500 Internal Server Error",
		}))
	}

	// Update fields that are allowed to be modified
	if requestData.EventName != "" {
		existingEvent.EventName = requestData.EventName
	}
	if requestData.EventDescription != "" {
		existingEvent.EventDescription = requestData.EventDescription
	}
	if requestData.CategoryName != "" {
		existingEvent.CategoryName = requestData.CategoryName
	}
	if requestData.CategoryId != "" {
		existingEvent.CategoryId = requestData.CategoryId
	}
	if requestData.EventType != "" {
		existingEvent.EventType = requestData.EventType
	}
	if requestData.EventMode != "" {
		existingEvent.EventMode = requestData.EventMode
	}
	if requestData.EventLocation != "" {
		existingEvent.EventLocation = requestData.EventLocation
	}
	if requestData.EventDate > 0 {
		existingEvent.EventDate = requestData.EventDate
	}
	if requestData.FlierImage != "" {
		existingEvent.FlierImage = requestData.FlierImage
	}
	if requestData.PaymentType != "" {
		existingEvent.PaymentType = requestData.PaymentType
	}
	if len(requestData.TicketComboDetails) > 0 {
		existingEvent.TicketComboDetails = requestData.TicketComboDetails
	}
	if requestData.Guidelines != "" {
		existingEvent.Guidelines = requestData.Guidelines
	}
	// if requestData.RegistrationLimit > 0 {
	// 	existingEvent.RegistrationLimit = requestData.RegistrationLimit
	// }

	// Update the updatedAt field
	existingEvent.UpdatedAt = time.Now().Unix()

	// Update the event in the database
	_, err = col.UpdateOne(ctx.Context(), bson.M{"uniqueId": eventID}, bson.M{
		"$set": existingEvent,
	})
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Failed to update event",
			Status:  "500 Internal Server Error",
		}))
	}

	// Return success response
	return ctx.JSON(commonutils.CreateSuccessResponse(&common_responses.SuccessResponse{
		Message: "Event updated successfully",
		Status:  "200 OK",
		Data: fiber.Map{
			"eventId": eventID,
		},
	}))
}

func DeleteEvent(ctx *fiber.Ctx) error {
	var requestData dbModel.Event

	// Parse request body
	err := json.Unmarshal(ctx.Body(), &requestData)
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error parsing request data",
			Status:  "400 Bad Request",
		}))
	}

	// Validate the event ID
	eventID := requestData.UniqueId
	if eventID == "" {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Event ID is required",
			Status:  "400 Bad Request",
		}))
	}

	// Connect to MongoDB
	db, col, err := mongoSetup.ConnectMongo("events")
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error connecting to MongoDB",
			Status:  "500 Internal Server Error",
		}))
	}
	defer db.Client().Disconnect(context.TODO())

	// Find the event
	var existingEvent dbModel.Event
	err = col.FindOne(ctx.Context(), bson.M{"uniqueId": eventID}).Decode(&existingEvent)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
				Message: "Event not found",
				Status:  "404 Not Found",
			}))
		}
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error fetching event",
			Status:  "500 Internal Server Error",
		}))
	}

	// Mark the event as inactive and update the updatedAt field
	existingEvent.Status = "inactive"
	existingEvent.UpdatedAt = time.Now().Unix()

	// Update the event's status in the database
	_, err = col.UpdateOne(ctx.Context(), bson.M{"uniqueId": eventID}, bson.M{
		"$set": bson.M{
			"status":    existingEvent.Status,
			"updatedAt": existingEvent.UpdatedAt,
		},
	})
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Failed to deactivate event",
			Status:  "500 Internal Server Error",
		}))
	}

	// Return success response
	return ctx.JSON(commonutils.CreateSuccessResponse(&common_responses.SuccessResponse{
		Message: "Event deactivated successfully",
		Status:  "200 OK",
		Data: fiber.Map{
			"eventId": eventID,
		},
	}))
}
