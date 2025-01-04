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

// CreateEvent handles the creation of an event.
func CreateEvent(ctx *fiber.Ctx) error {
	// Parse request body
	var requestData dbModel.Event
	if err := ctx.BodyParser(&requestData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing request data",
			"status":  "400 Bad Request",
		})
	}
	fmt.Println(requestData)

	// Validate required fields
	requestData.EventName = strings.TrimSpace(requestData.EventName)
	if requestData.EventName == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Event name is required",
			"status":  "400 Bad Request",
		})
	}

	requestData.EventDescription = strings.TrimSpace(requestData.EventDescription)
	if requestData.EventDescription == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Event description is required",
			"status":  "400 Bad Request",
		})
	}

	if requestData.ParticipationCapacity <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Participant limit must be greater than 0",
			"status":  "400 Bad Request",
		})
	}

	currentTime := time.Now().Unix()
	if requestData.EventDate < currentTime {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Date of event cannot be in the past",
			"status":  "400 Bad Request",
		})
	}

	requestData.EventMode = strings.TrimSpace(requestData.EventMode)
	if requestData.EventMode != "online" && requestData.EventMode != "offline" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Event mode must be 'online' or 'offline'",
			"status":  "400 Bad Request",
		})
	}

	requestData.PaymentType = strings.TrimSpace(requestData.PaymentType)
	if requestData.PaymentType != "free" && requestData.PaymentType != "paid" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Payment type must be 'free' or 'paid'",
			"status":  "400 Bad Request",
		})
	}

	if requestData.PaymentType == "paid" && requestData.RegistrationAmount <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Registration amount must be greater than 0 for paid events",
			"status":  "400 Bad Request",
		})
	}

	// Validate registration form fields
	// for _, field := range requestData.RegistrationForm {
	// 	field.Label = strings.TrimSpace(field.Label)
	// 	field.Type = strings.TrimSpace(field.Type)

	// 	if field.Label == "" {
	// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 			"message": "Form field label is required",
	// 			"status":  "400 Bad Request",
	// 		})
	// 	}

	// 	if field.Type != "textarea" && field.Type != "input" && field.Type != "select" {
	// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 			"message": "Unsupported form field type: " + field.Type,
	// 			"status":  "400 Bad Request",
	// 		})
	// 	}

	// 	if field.Type == "select" && len(field.Data) == 0 {
	// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 			"message": "Select field must have options",
	// 			"status":  "400 Bad Request",
	// 		})
	// 	}
	// }

	// Generate unique ID and timestamps
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

	// Connect to MongoDB
	db, col, err := mongoSetup.ConnectMongo("events")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error connecting to MongoDB",
			"status":  "500 Internal Server Error",
		})
	}
	defer db.Client().Disconnect(context.Background())

	// Insert into MongoDB
	_, err = col.InsertOne(ctx.Context(), requestData)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to insert event",
			"status":  "500 Internal Server Error",
		})
	}

	// Return success response
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Event created successfully",
		"status":  "201 Created",
		"data": fiber.Map{
			"eventId": requestData.UniqueId,
		},
	})
}

func EditEvent(ctx *fiber.Ctx) error {
	// loginDetails := ctx.Locals("userData").(common_responses.LoginDetails)
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

	// Find the event in the database
	db, col, err := mongoSetup.ConnectMongo("events")
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error connecting to MongoDB",
			Status:  "500 Internal Server Error",
		}))
	}
	defer db.Client().Disconnect(context.TODO())

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

	// Update event fields
	if requestData.EventName != "" {
		existingEvent.EventName = requestData.EventName
	}
	if requestData.EventDescription != "" {
		existingEvent.EventDescription = requestData.EventDescription
	}
	if requestData.ParticipationCapacity > 0 {
		existingEvent.ParticipationCapacity = requestData.ParticipationCapacity
	}
	if requestData.EventDate > 0 {
		existingEvent.EventDate = requestData.EventDate
	}
	if requestData.EventMode != "" {
		existingEvent.EventMode = requestData.EventMode
	}
	if requestData.PaymentType != "" {
		existingEvent.PaymentType = requestData.PaymentType
	}
	if requestData.RegistrationAmount > 0 {
		existingEvent.RegistrationAmount = requestData.RegistrationAmount
	}

	// Update the updatedAt field to the current timestamp
	existingEvent.UpdatedAt = time.Now().Unix()

	// Update the event in the database
	_, err = col.UpdateOne(ctx.Context(), bson.M{"uniqueId": eventID}, bson.M{
		"$set": bson.M{
			"eventName":             existingEvent.EventName,
			"eventDescription":      existingEvent.EventDescription,
			"participationCapacity": existingEvent.ParticipationCapacity,
			"eventDate":             existingEvent.EventDate,
			"eventMode":             existingEvent.EventMode,
			"paymentType":           existingEvent.PaymentType,
			"registrationAmount":    existingEvent.RegistrationAmount,
			"updatedAt":             existingEvent.UpdatedAt,
		},
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

	// Mark the event as inactive and update the UpdatedAt field
	existingEvent.Status = "inactive"
	existingEvent.UpdatedAt = time.Now().Unix() // Set the updated timestamp

	// Update the event's status and UpdatedAt in the database
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
