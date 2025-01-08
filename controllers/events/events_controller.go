package eventPanel

import (
	"context"
	mongoSetup "em_backend/configs/mongo"
	commonutils "em_backend/library/common"
	dbModel "em_backend/models/db"
	common_responses "em_backend/responses/common"
	event_response "em_backend/responses/event"
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

func GetEventByID(ctx *fiber.Ctx) error {
	// Parse request body
	var requestData dbModel.Event
	err := json.Unmarshal(ctx.Body(), &requestData)
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error parsing request data",
			Status:  "400 Bad Request",
		}))
	}

	// Get the event ID from the parsed request data
	eventId := requestData.UniqueId
	if eventId == "" {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Event ID is required",
			Status:  "400 Bad Request",
		}))
	}

	// Define filter and projection for the query
	filter := bson.M{"uniqueId": eventId}
	projection := bson.M{} // Add specific fields if needed; leave empty for all fields

	// Use the `FindOneDoc` function to query the database
	result, err := mongoSetup.FindOneDoc("events", filter, projection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
				Message: "Event not found",
				Status:  "404 Not Found",
			}))
		}
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: fmt.Sprintf("Error fetching event: %v", err),
			Status:  "500 Internal Server Error",
		}))
	}

	// Decode the result into the Event struct
	var event dbModel.Event
	if decodeErr := result.Decode(&event); decodeErr != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error decoding event data",
			Status:  "500 Internal Server Error",
		}))
	}

	// Return the event as a success response
	return ctx.JSON(commonutils.CreateSuccessResponse(&common_responses.SuccessResponse{
		Message: "Event fetched successfully",
		Status:  "200 OK",
		Data:    event,
	}))
}

func GetRegistrationForm(ctx *fiber.Ctx) error {
	// Get the event ID from the query parameter
	var requestData dbModel.Event

	// Parse request body
	err := json.Unmarshal(ctx.Body(), &requestData)
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error parsing request data",
			Status:  "400 Bad Request",
		}))
	}
	fmt.Println("==", requestData)

	if requestData.UniqueId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "eventId is required",
			"status":  "400 Bad Request",
		})
	}

	// MongoDB Query
	filter := bson.M{"uniqueId": requestData.UniqueId}

	// Use the `FindOneDoc` function to query the database
	result, err := mongoSetup.FindOneDoc("events", filter, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
				Message: "Event not found",
				Status:  "404 Not Found",
			}))
		}
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: fmt.Sprintf("Error fetching event: %v", err),
			Status:  "500 Internal Server Error",
		}))
	}

	// Decode the result into the Event struct
	var event dbModel.Event
	if decodeErr := result.Decode(&event); decodeErr != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error decoding event data",
			Status:  "500 Internal Server Error",
		}))
	}
	fmt.Println("filt", filter)
	formInfo, err := mongoSetup.FindOneDoc("registerForm", bson.M{"eventId": requestData.UniqueId}, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Registration form not found",
				"status":  "404 Not Found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching registration form",
			"status":  "500 Internal Server Error",
		})
	}
	// Decode the result into the Event struct
	var formDetails dbModel.RegistrationForm
	if decodeErr := formInfo.Decode(&formDetails); decodeErr != nil {
		fmt.Println(decodeErr)
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error decoding register form data",
			Status:  "500 Internal Server Error",
		}))
	}
	// Success Response
	fmt.Println("for", formDetails)
	var responseData event_response.EventRegisterFormInfoResp
	responseData.FormFields = formDetails
	responseData.EventName = event.EventName
	fmt.Println(responseData)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Registration form fetched successfully",
		"status":  "200 OK",
		"data":    responseData,
	})
}

func RegisterEvent(ctx *fiber.Ctx) error {
	var requestData map[string]interface{} // Dynamic structure for request data

	// Parse request body
	err := json.Unmarshal(ctx.Body(), &requestData)
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error parsing request data",
			Status:  "400 Bad Request",
		}))
	}

	// Validate Event ID
	eventID, ok := requestData["eventId"].(string)
	if !ok || eventID == "" {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Event ID is required",
			Status:  "400 Bad Request",
		}))
	}

	// Add metadata fields
	requestData["createdAt"] = time.Now().Unix()
	requestData["updatedAt"] = time.Now().Unix()

	// Save the data as-is into the registrations collection
	_, registrationsCol, err := mongoSetup.ConnectMongo("registrations")
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error connecting to Registrations collection",
			Status:  "500 Internal Server Error",
		}))
	}
	defer registrationsCol.Database().Client().Disconnect(context.TODO())

	_, err = registrationsCol.InsertOne(ctx.Context(), requestData)
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
			"eventId": eventID,
		},
	}))
}

func GetRegistrationDetails(ctx *fiber.Ctx) error {
	// Get eventId from query parameters
	eventID := ctx.Query("eventId")
	if eventID == "" {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Event ID is required",
			Status:  "400 Bad Request",
		}))
	}

	// Connect to MongoDB
	_, registrationsCol, err := mongoSetup.ConnectMongo("registrations")
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error connecting to Registrations collection",
			Status:  "500 Internal Server Error",
		}))
	}
	defer registrationsCol.Database().Client().Disconnect(context.TODO())

	// Query the database for registrations of the given event ID
	filter := bson.M{"eventId": eventID}
	cursor, err := registrationsCol.Find(ctx.Context(), filter)
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error fetching registration details",
			Status:  "500 Internal Server Error",
		}))
	}
	defer cursor.Close(ctx.Context())

	// Process the results
	var registrations []map[string]interface{}
	err = cursor.All(ctx.Context(), &registrations)
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error processing registration details",
			Status:  "500 Internal Server Error",
		}))
	}

	// Convert registrations into the desired structure
	var result []map[string]interface{}
	for _, registration := range registrations {
		entry := []map[string]interface{}{}
		for key, value := range registration {
			// Skip MongoDB metadata fields
			if key == "_id" || key == "createdAt" || key == "updatedAt" {
				continue
			}
			entry = append(entry, map[string]interface{}{
				"label": key,
				"data":  value,
			})
		}
		result = append(result, map[string]interface{}{
			"registration": entry,
		})
	}

	// Return the response
	return ctx.JSON(commonutils.CreateSuccessResponse(&common_responses.SuccessResponse{
		Message: "Registration details fetched successfully",
		Status:  "200 OK",
		Data:    result,
	}))
}

// func RegisterEvent(ctx *fiber.Ctx) error {
// 	var requestData map[string]interface{} // Dynamic structure for request data

// 	// Parse request body
// 	err := json.Unmarshal(ctx.Body(), &requestData)
// 	if err != nil {
// 		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
// 			Message: "Error parsing request data",
// 			Status:  "400 Bad Request",
// 		}))
// 	}

// 	// Validate Event ID
// 	eventID, ok := requestData["eventId"].(string)
// 	if !ok || eventID == "" {
// 		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
// 			Message: "Event ID is required",
// 			Status:  "400 Bad Request",
// 		}))
// 	}

// 	// Connect to MongoDB
// 	db, eventsCol, err := mongoSetup.ConnectMongo("events")
// 	if err != nil {
// 		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
// 			Message: "Error connecting to MongoDB",
// 			Status:  "500 Internal Server Error",
// 		}))
// 	}
// 	defer db.Client().Disconnect(context.TODO())

// 	// Fetch event details
// 	var event dbModel.Event
// 	err = eventsCol.FindOne(ctx.Context(), bson.M{"uniqueId": eventID, "status": "active"}).Decode(&event)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
// 				Message: "Event not found or not active",
// 				Status:  "404 Not Found",
// 			}))
// 		}
// 		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
// 			Message: "Error fetching event details",
// 			Status:  "500 Internal Server Error",
// 		}))
// 	}

// 	// Fetch registration form details
// 	filter := bson.M{"eventId": eventID}
// 	formDoc, err := mongoSetup.FindOneDoc("registerForm", filter, bson.M{})
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
// 				Message: "Registration form not found for the event",
// 				Status:  "404 Not Found",
// 			}))
// 		}
// 		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
// 			Message: "Error fetching registration form",
// 			Status:  "500 Internal Server Error",
// 		}))
// 	}

// 	var form dbModel.RegistrationForm
// 	if decodeErr := formDoc.Decode(&form); decodeErr != nil {
// 		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
// 			Message: "Error decoding registration form",
// 			Status:  "500 Internal Server Error",
// 		}))
// 	}

// 	// Validate requestData against form fields
// 	for _, field := range form.RegistrationFormFields {
// 		fieldValue, exists := requestData[field.Label]
// 		if !exists {
// 			return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
// 				Message: fmt.Sprintf("Missing required field: %s", field.Label),
// 				Status:  "400 Bad Request",
// 			}))
// 		}

// 		switch field.Type {
// 		case "text":
// 			if _, ok := fieldValue.(string); !ok {
// 				return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
// 					Message: fmt.Sprintf("Invalid data type for field: %s, expected text", field.Label),
// 					Status:  "400 Bad Request",
// 				}))
// 			}
// 		case "arrayOfStrings":
// 			if _, ok := fieldValue.([]string); !ok {
// 				return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
// 					Message: fmt.Sprintf("Invalid data type for field: %s, expected array of strings", field.Label),
// 					Status:  "400 Bad Request",
// 				}))
// 			}
// 		default:
// 			return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
// 				Message: fmt.Sprintf("Unsupported field type: %s", field.Type),
// 				Status:  "400 Bad Request",
// 			}))
// 		}
// 	}

// 	// Add metadata and prepare for insertion
// 	requestData["createdAt"] = time.Now().Unix()
// 	requestData["updatedAt"] = time.Now().Unix()

// 	// Save the validated data into the registrations collection
// 	_, registrationsCol, err := mongoSetup.ConnectMongo("registrations")
// 	if err != nil {
// 		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
// 			Message: "Error connecting to Registrations collection",
// 			Status:  "500 Internal Server Error",
// 		}))
// 	}

// 	_, err = registrationsCol.InsertOne(ctx.Context(), requestData)
// 	if err != nil {
// 		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
// 			Message: "Failed to register for the event",
// 			Status:  "500 Internal Server Error",
// 		}))
// 	}

// 	// Return success response
// 	return ctx.JSON(commonutils.CreateSuccessResponse(&common_responses.SuccessResponse{
// 		Message: "Event registered successfully",
// 		Status:  "200 OK",
// 		Data: fiber.Map{
// 			"eventId": eventID,
// 		},
// 	}))
// }
