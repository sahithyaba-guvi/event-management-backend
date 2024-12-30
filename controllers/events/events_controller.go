package eventPanel

import (
	"context"
	mongoSetup "em_backend/configs/mongo"
	commonutils "em_backend/library/common"
	dbModel "em_backend/models/db"
	common_responses "em_backend/responses/common"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
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
