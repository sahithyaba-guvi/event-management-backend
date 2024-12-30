package adminpanel

import (
	"context"
	mongoSetup "em_backend/configs/mongo"
	"encoding/json"
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
	loginDetails := ctx.Locals("userData").(common_responses.LoginDetails)
	var requestData dbModel.Event

	// Parse request body
	err := json.Unmarshal(ctx.Body(), &requestData)
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error parsing request data",
			Status:  "400 Bad Request",
		}))
	}

	// Validate fields separately
	requestData.EventName = strings.TrimSpace(requestData.EventName)
	if requestData.EventName == "" {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Event name is required",
			Status:  "400 Bad Request",
		}))
	}

	requestData.EventDescription = strings.TrimSpace(requestData.EventDescription)
	if requestData.EventDescription == "" {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Event description is required",
			Status:  "400 Bad Request",
		}))
	}

	if requestData.ParticipantLimit <= 0 {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Participant limit must be greater than 0",
			Status:  "400 Bad Request",
		}))
	}

	currentTime := time.Now().Unix()
	if requestData.DateOfEvent < currentTime {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Date of event cannot be in the past",
			Status:  "400 Bad Request",
		}))
	}

	requestData.OnlineOffline = strings.TrimSpace(requestData.OnlineOffline)
	if requestData.OnlineOffline != "online" && requestData.OnlineOffline != "offline" {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "OnlineOffline must be either 'Online' or 'Offline'",
			Status:  "400 Bad Request",
		}))
	}

	requestData.FreePaid = strings.TrimSpace(requestData.FreePaid)
	if requestData.FreePaid != "free" && requestData.FreePaid != "paid" {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "FreePaid must be either 'Free' or 'Paid'",
			Status:  "400 Bad Request",
		}))
	}

	if requestData.FreePaid == "paid" && requestData.Amount <= 0 {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Amount must be greater than 0 for paid events",
			Status:  "400 Bad Request",
		}))
	}

	requestData.AdminHash = loginDetails.Hash

	// Assign a unique ID and timestamps
	eventID, err := uuid.NewRandom()
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error generating event ID",
			Status:  "500 Internal Server Error",
		}))
	}
	requestData.UniqueId = eventID.String()
	requestData.CreatedAt = time.Now().Unix()

	// Connect to MongoDB
	db, col, err := mongoSetup.ConnectMongo("events")
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error connecting to MongoDB",
			Status:  "500 Internal Server Error",
		}))
	}
	defer db.Client().Disconnect(context.TODO())

	// Insert the event into MongoDB
	_, err = col.InsertOne(ctx.Context(), requestData)
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Failed to insert event",
			Status:  "500 Internal Server Error",
		}))
	}

	// Return success response
	return ctx.JSON(commonutils.CreateSuccessResponse(&common_responses.SuccessResponse{
		Message: "Event created successfully",
		Status:  "201 Created",
		Data: fiber.Map{
			"eventId": requestData.UniqueId,
		},
	}))
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
	eventID := ctx.Params("eventId")
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
	if requestData.ParticipantLimit > 0 {
		existingEvent.ParticipantLimit = requestData.ParticipantLimit
	}
	if requestData.DateOfEvent > 0 {
		existingEvent.DateOfEvent = requestData.DateOfEvent
	}
	if requestData.OnlineOffline != "" {
		existingEvent.OnlineOffline = requestData.OnlineOffline
	}
	if requestData.FreePaid != "" {
		existingEvent.FreePaid = requestData.FreePaid
	}
	if requestData.Amount > 0 {
		existingEvent.Amount = requestData.Amount
	}

	// Update the updatedAt field to the current timestamp
	existingEvent.UpdatedAt = time.Now().Unix()

	// Update the event in the database
	_, err = col.UpdateOne(ctx.Context(), bson.M{"uniqueId": eventID}, bson.M{
		"$set": bson.M{
			"eventName":        existingEvent.EventName,
			"eventDescription": existingEvent.EventDescription,
			"participantLimit": existingEvent.ParticipantLimit,
			"dateOfEvent":      existingEvent.DateOfEvent,
			"onlineOffline":    existingEvent.OnlineOffline,
			"freePaid":         existingEvent.FreePaid,
			"amount":           existingEvent.Amount,
			"updatedAt":        existingEvent.UpdatedAt,
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
	// Get event ID from URL params
	eventID := ctx.Params("eventId")
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
