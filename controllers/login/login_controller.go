package loginPanel

import (
	mongoSetup "em_backend/configs/mongo"
	common "em_backend/library/common"
	loginCommon "em_backend/library/login"
	"encoding/json"
	"strings"

	dbModel "em_backend/models/db"
	loginModel "em_backend/models/login"
	common_responses "em_backend/responses/common"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Login(ctx *fiber.Ctx) error {
	var loginRequest loginModel.LoginReq
	fmt.Println("it came to login")
	// Parse the login request
	if err := ctx.BodyParser(&loginRequest); err != nil {
		return ctx.JSON(common.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error parsing request data",
			Status:  "400 Bad Request",
		}))
	}

	// Hash the provided password (assuming bcrypt is used for hashing)
	hashedPassword := loginRequest.Password

	// Find user by email
	filter := bson.M{"email": loginRequest.Email}
	projection := bson.M{"_id": 1, "password": 1, "userHash": 1, "userName": 1, "email": 1}

	result, err := mongoSetup.FindOneDoc("userData", filter, projection)
	if err != nil {
		fmt.Println("Error querying database:", err)
		return ctx.JSON(common.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error querying database",
			Status:  "500 Internal Server Error",
		}))
	}
	var user dbModel.UserData
	// Check if the user exists and the password matches
	if result.Decode(&user) != nil {
		return ctx.JSON(common.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Invalid email or password",
			Status:  "401 Unauthorized",
		}))
	}
	fmt.Println(user)

	// Compare hashed password with the stored password (bcrypt check)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(hashedPassword))
	if err != nil {
		// Passwords don't match
		return ctx.JSON(common.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Invalid email or password",
			Status:  "401 Unauthorized",
		}))
	}

	// Find user by email

	// Generate Auth Token
	authToken, err := loginCommon.GenerateAuthToken(user.UserHash)
	if err != nil {
		return ctx.JSON(common.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error generating auth token",
			Status:  "500 Internal Server Error",
		}))
	}

	var userInfo common_responses.LoginDetails
	userInfo.UserName = user.UserName
	userInfo.Email = user.Email
	userInfo.Hash = user.UserHash
	userInfo.AuthToken = authToken
	userInfo.IsAdmin = common.CheckAdmin(user.Email)

	// Store the token in Redis
	err = loginCommon.StoreAuthTokenInRedis(userInfo)
	if err != nil {
		return ctx.JSON(common.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error storing auth token in Redis",
			Status:  "500 Internal Server Error",
		}))
	}

	// Return success response (user is authenticated)
	return ctx.JSON(common.CreateSuccessResponse(&common_responses.SuccessResponse{
		Message: "Login successful",
		Status:  "200 OK",
		Data: fiber.Map{
			"authToken": authToken,
			"name":      user.UserName,
			"email":     user.Email,
			"isAdmin":   userInfo.IsAdmin,
		},
	}))
}

func Register(ctx *fiber.Ctx) error {
	var registrationRequest loginModel.RegisterReq

	// Parse the incoming request body
	if err := ctx.BodyParser(&registrationRequest); err != nil {
		return ctx.JSON(common.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Invalid request data",
			Status:  "400 Bad Request",
		}))
	}

	// Validate required fields
	if registrationRequest.UserName == "" || registrationRequest.Email == "" || registrationRequest.Password == "" {
		return ctx.JSON(common.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "All fields are required",
			Status:  "400 Bad Request",
		}))
	}

	// Check if email is already registered
	emailExists, err := loginCommon.IsEmailRegistered(registrationRequest.Email)
	if err != nil {
		return ctx.JSON(common.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error checking email existence",
			Status:  "500 Internal Server Error",
		}))
	}

	if emailExists {
		return ctx.JSON(common.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Email already registered",
			Status:  "409 Conflict",
		}))
	}

	// Hash the password
	hashedPassword, err := loginCommon.HashPassword(registrationRequest.Password)
	if err != nil {
		return ctx.JSON(common.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error securing the password",
			Status:  "500 Internal Server Error",
		}))
	}

	// Register the new user

	// Generate a unique hash for the user
	userHash := loginCommon.GenerateUserHash(registrationRequest.Email)

	// Register a new user
	var user dbModel.UserData
	user.Email = registrationRequest.Email
	user.UserName = registrationRequest.UserName
	user.Password = hashedPassword
	user.UserHash = userHash
	user.CreatedAt = time.Now().Unix()
	_, err = mongoSetup.InsertOneDoc("userData", user)
	if err != nil {
		return ctx.JSON(common.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Failed to register user",
			Status:  "500 Internal Server Error",
		}))
	}

	// Generate Auth Token
	authToken, err := loginCommon.GenerateAuthToken(userHash)
	if err != nil {
		return ctx.JSON(common.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error generating auth token",
			Status:  "500 Internal Server Error",
		}))
	}
	var userInfo common_responses.LoginDetails
	userInfo.UserName = registrationRequest.UserName
	userInfo.Email = registrationRequest.Email
	userInfo.Hash = userHash
	userInfo.AuthToken = authToken

	// Store the token in Redis
	err = loginCommon.StoreAuthTokenInRedis(userInfo)
	if err != nil {
		return ctx.JSON(common.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error storing auth token in Redis",
			Status:  "500 Internal Server Error",
		}))
	}

	// Return success response (user is registered and authenticated)
	return ctx.JSON(common.CreateSuccessResponse(&common_responses.SuccessResponse{
		Message: "User registered successfully",
		Status:  "201 Created",
		Data: fiber.Map{
			"authToken": authToken,
			"date":      userInfo,
			"isAdmin":   common.CheckAdmin(userInfo.Email),
		},
	}))
}

func RegisterForgotPasswordEmail(ctx *fiber.Ctx) error {
	var res common_responses.Response
	var email loginModel.RegisterNewUserOrResetPasswordReq
	requestBody := ctx.Body()
	if len(requestBody) != 0 {
		err := json.Unmarshal(requestBody, &email)
		if err != nil {
			fmt.Println("login requestBody unmarshal error:", err)
			res.Access = false
			res.Message = "Invalid/Empty Request body"
			res.Status = "400 Bad Request"
			return ctx.JSON(res)
		}
	} else {
		res.Access = false
		res.Message = "Invalid/Empty Request body"
		res.Status = "400 Bad Request"
		return ctx.JSON(res)
	}

	emailVal := strings.TrimSpace(strings.ToLower(email.Email))
	emailExists, err := loginCommon.IsEmailRegistered(email.Email)
	if err != nil {
		return ctx.JSON(common.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error checking email existence",
			Status:  "500 Internal Server Error",
		}))
	}

	if emailExists {
		return ctx.JSON(common.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Email already registered",
			Status:  "409 Conflict",
		}))
	}

	isEmailSent, err := loginCommon.SendRegisterOrForgotPasswordEmail(emailVal, "register")
	if !isEmailSent {
		res.Access = false
		res.Message = "Cannot send email right now"
		res.Status = "409 Conflict"
		return ctx.JSON(res)
	}
	res.Access = true
	res.Message = "Reset password mail sent succesfully"
	res.Status = "200 OK"
	return ctx.JSON(res)
}

func VerifyOTPForgotPassword(ctx *fiber.Ctx) error {
	var request loginModel.ForgotPasswordOTPRequest
	var response common_responses.Response

	requestBody := ctx.Body()

	if len(requestBody) != 0 {
		err := json.Unmarshal(requestBody, &request)
		if err != nil {
			fmt.Println("login requestBody unmarshal error:", err)
			response.Access = false
			response.Message = "Invalid/Empty Request body"
			response.Status = "400 Bad Request"
			return ctx.JSON(response)
		}
	} else {
		response.Access = false
		response.Message = "Invalid/Empty Request body"
		response.Status = "400 Bad Request"
		return ctx.JSON(response)
	}

	isOTPCorrect, err := loginCommon.VerifyForgotPasswordOTP(request.Email, request.OTP)

	if err != nil {
		fmt.Println("Cannot verify password right now", err)
		response.Access = false
		response.Message = "Incorrect OTP"
		response.Status = "409 Conflict"
		return ctx.JSON(response)
	}

	if !isOTPCorrect {
		response.Access = false
		response.Message = "Incorrect OTP"
		response.Status = "409 Conflict"
		return ctx.JSON(response)
	}

	response.Access = true
	response.Message = "Correct OTP"
	response.Status = "200 OK"
	return ctx.JSON(response)

}
