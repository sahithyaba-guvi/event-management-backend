package middleware

import (
	common "em_backend/library/common"
	loginFunc "em_backend/library/login"
	commonModel "em_backend/models/common"
	common_resp "em_backend/responses/common"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func AuthenticationMiddlewareForAdmin(ctx *fiber.Ctx) error {
	// Declare a variable to store the auth token from the request body
	var data commonModel.Authtoken
	fmt.Println("it came rhere")
	// Read the request body
	requestBody := ctx.Body()

	// Check if the request body is not empty
	if len(requestBody) == 0 {
		// Return a 400 Bad Request response if the body is missing or empty
		response := common_resp.FailureResponse{
			Access:  false,
			Status:  "400 Bad Request",
			Message: "Empty/Missing request body",
		}
		return ctx.JSON(response)
	}

	// Unmarshal the request body into the `data` struct
	err := json.Unmarshal(requestBody, &data)
	if err != nil {
		// Log and return an error if unmarshalling fails
		fmt.Println("Error Unmarshaling the request body:", err)
		response := common_resp.FailureResponse{
			Access:  false,
			Status:  "400 Bad Request",
			Message: "Invalid request format",
		}
		return ctx.JSON(response)
	}

	// Check if the auth token is present in the request
	if data.AuthToken == "" {
		// Return a 400 Bad Request response if the token is missing
		response := common_resp.FailureResponse{
			Access:  false,
			Status:  "400 Bad Request",
			Message: "Missing AuthToken",
		}
		return ctx.JSON(response)
	}

	// Check the session validity using the auth token
	session := loginFunc.CheckSession(data.AuthToken)

	// If session is invalid, return a 401 Unauthorized response
	if !session.Access {
		response := common_resp.FailureResponse{
			Access:  false,
			Status:  "401 Unauthorized",
			Message: "Session invalid",
		}
		return ctx.JSON(response)
	}
	if !common.CheckAdmin(session.UserInfo.Email) {
		response := common_resp.FailureResponse{
			Access:  false,
			Status:  "401 Unauthorized",
			Message: "Not an admin",
		}
		return ctx.JSON(response)
	}
	ctx.Locals("userData", session.UserInfo)
	fmt.Println("evetin ok ------------")
	// If the session is valid, proceed to the next middleware/handler
	return ctx.Next()
}

func AuthenticationMiddleware(ctx *fiber.Ctx) error {
	// Declare a variable to store the auth token from the request body
	var data commonModel.Authtoken
	fmt.Println("it came rhere")
	// Read the request body
	requestBody := ctx.Body()

	// Check if the request body is not empty
	if len(requestBody) == 0 {
		// Return a 400 Bad Request response if the body is missing or empty
		response := common_resp.FailureResponse{
			Access:  false,
			Status:  "400 Bad Request",
			Message: "Empty/Missing request body",
		}
		return ctx.JSON(response)
	}

	// Unmarshal the request body into the `data` struct
	err := json.Unmarshal(requestBody, &data)
	if err != nil {
		// Log and return an error if unmarshalling fails
		fmt.Println("Error Unmarshaling the request body:", err)
		response := common_resp.FailureResponse{
			Access:  false,
			Status:  "400 Bad Request",
			Message: "Invalid request format",
		}
		return ctx.JSON(response)
	}

	// Check if the auth token is present in the request
	if data.AuthToken == "" {
		// Return a 400 Bad Request response if the token is missing
		response := common_resp.FailureResponse{
			Access:  false,
			Status:  "400 Bad Request",
			Message: "Missing AuthToken",
		}
		return ctx.JSON(response)
	}

	// Check the session validity using the auth token
	session := loginFunc.CheckSession(data.AuthToken)

	// If session is invalid, return a 401 Unauthorized response
	if !session.Access {
		response := common_resp.FailureResponse{
			Access:  false,
			Status:  "401 Unauthorized",
			Message: "Session invalid",
		}
		return ctx.JSON(response)
	}
	session.UserInfo.IsAdmin = common.CheckAdmin(session.UserInfo.Email)
	ctx.Locals("userData", session.UserInfo)
	fmt.Println("evetin ok ------------")
	// If the session is valid, proceed to the next middleware/handler
	return ctx.Next()
}
