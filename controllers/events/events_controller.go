package eventPanel

import (
	"context"
	mongoSetup "em_backend/configs/mongo"
	commonutils "em_backend/library/common"
	dbModel "em_backend/models/db"
	common_responses "em_backend/responses/common"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllEvents(ctx *fiber.Ctx) error {
	// Connect to MongoDB
	db, col, err := mongoSetup.ConnectMongo("events")
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error connecting to MongoDB",
			Status:  "500 Internal Server Error",
		}))
	}
	defer db.Client().Disconnect(context.TODO())

	// Define filter to fetch only active events
	filter := bson.M{"status": "active"}

	// Fetch events from the collection
	cursor, err := col.Find(ctx.Context(), filter, options.Find())
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error fetching events",
			Status:  "500 Internal Server Error",
		}))
	}
	defer cursor.Close(ctx.Context())

	var events []dbModel.Event
	if err = cursor.All(ctx.Context(), &events); err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error processing events",
			Status:  "500 Internal Server Error",
		}))
	}
	fmt.Println("==eve", events)
	// Return the events as a success response
	return ctx.JSON(commonutils.CreateSuccessResponse(&common_responses.SuccessResponse{
		Message: "Events fetched successfully",
		Status:  "200 OK",
		Data:    events,
	}))
}

func RegisterEvent(ctx *fiber.Ctx) error {
	var requestData dbModel.Registration

	// Parse request body
	err := json.Unmarshal(ctx.Body(), &requestData)
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error parsing request data",
			Status:  "400 Bad Request",
		}))
	}

	// Validate event ID
	eventID := requestData.EventId
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

	// Check if event exists and is active
	var existingEvent dbModel.Event
	err = col.FindOne(ctx.Context(), bson.M{"uniqueId": eventID, "status": "active"}).Decode(&existingEvent)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
				Message: "Event not found or not active",
				Status:  "404 Not Found",
			}))
		}
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error fetching event details",
			Status:  "500 Internal Server Error",
		}))
	}

	// Prepare the registration data
	registrationData := dbModel.Registration{
		RegisterId:        requestData.RegisterId,
		EventId:           requestData.EventId,
		TeamSize:          requestData.TeamSize,
		TeamMemberDetails: requestData.TeamMemberDetails,
		PrimaryEmailId:    requestData.PrimaryEmailId,
		CreatedAt:         time.Now().Unix(),
		UpdatedAt:         time.Now().Unix(),
	}

	// Save the registration data into the registrations collection
	_, registrationsCol, err := mongoSetup.ConnectMongo("registrations")
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error connecting to Registrations collection",
			Status:  "500 Internal Server Error",
		}))
	}

	_, err = registrationsCol.InsertOne(ctx.Context(), registrationData)
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Failed to register for the event",
			Status:  "500 Internal Server Error",
		}))
	}

	// Return success response
	return ctx.JSON(commonutils.CreateSuccessResponse(&common_responses.SuccessResponse{
		Message: "Event registered successfully",
		Status:  "200 OK",
		Data: fiber.Map{
			"registerId": requestData.RegisterId,
			"eventId":    requestData.EventId,
		},
	}))
}
