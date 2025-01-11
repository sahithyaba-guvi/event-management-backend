package paymentPanel

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	mongoSetup "em_backend/configs/mongo"

	commonutils "em_backend/library/common"
	paymentModel "em_backend/models/payment"
	common_responses "em_backend/responses/common"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

// createOrderHandler handles order creation
func CreateOrderHandler(ctx *fiber.Ctx) error {
	// Create order request
	var orderReq paymentModel.OrderRequest
	err := json.Unmarshal(ctx.Body(), &orderReq)
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error parsing request data",
			Status:  "400 Bad Request",
		}))
	}
	// Call CreateOrder to create the Razorpay order
	order, err := CreateOrder(&orderReq)
	if err != nil {
		log.Printf("Failed to create order: %v", err)
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Failed to create Razorpay order",
			Status:  "500 Internal Server Error",
		}))
	}

	// Prepare order details for insertion
	orderDetails := paymentModel.OrderDetails{
		OrderID:  order.ID,
		Amount:   order.Amount,
		Currency: order.Currency,
		Receipt:  order.Receipt,
		Status:   order.Status,
		RegID:    orderReq.RegID,
	}

	// Connect to the MongoDB collection
	_, orderDetailsCol, err := mongoSetup.ConnectMongo("orderDetails")
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error connecting to OrderDetails collection",
			Status:  "500 Internal Server Error",
		}))
	}
	defer orderDetailsCol.Database().Client().Disconnect(context.TODO())

	// Insert order details into the MongoDB collection
	_, err = orderDetailsCol.InsertOne(ctx.Context(), orderDetails)
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Failed to save order details",
			Status:  "500 Internal Server Error",
		}))
	}

	// Return the success response
	return ctx.JSON(commonutils.CreateSuccessResponse(&common_responses.SuccessResponse{
		Message: "Order created and details saved successfully",
		Status:  "200 OK",
		Data:    order,
	}))
}

// createOrder creates a new Razorpay order
func CreateOrder(orderReq *paymentModel.OrderRequest) (*paymentModel.OrderResponse, error) {
	// Load configuration values
	apiSecret := commonutils.LoadEnv("RAZORPAY_API_SECRET")
	apiKey := commonutils.LoadEnv("RAZORPAY_API_KEY")
	baseURL := commonutils.LoadEnv("RAZORPAY_BASE_URL")

	// Initialize Resty client
	client := resty.New()

	// Perform the POST request
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(apiKey, apiSecret).
		SetBody(orderReq).
		Post(baseURL + "/orders")

	if err != nil {
		log.Printf("Error making API request: %v", err)
		return nil, err
	}

	// Check for non-200 HTTP response
	if resp.StatusCode() != 200 && resp.StatusCode() != 201 {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	// Parse the response body
	var orderResp paymentModel.OrderResponse
	if err := json.Unmarshal(resp.Body(), &orderResp); err != nil {
		log.Printf("Error parsing response body: %v", err)
		return nil, err
	}

	return &orderResp, nil
}

// verifyPaymentHandler verifies the Razorpay payment signature
func VerifyPaymentHandler(ctx *fiber.Ctx) error {
	// Parse the request body
	verification := new(paymentModel.PaymentVerification)
	if err := ctx.BodyParser(verification); err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Invalid request body",
			Status:  "400 Bad Request",
		}))
	}

	// Verify the payment signature
	if !verifySignature(verification) {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Invalid payment signature",
			Status:  "400 Bad Request",
		}))
	}

	// Prepare payment details
	paymentDetails := paymentModel.PaymentDetails{
		OrderID:   verification.OrderID,
		PaymentID: verification.PaymentID,
		Status:    "Verified",
	}

	// Connect to the MongoDB collection
	_, paymentDetailsCol, err := mongoSetup.ConnectMongo("paymentDetails")
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Error connecting to PaymentDetails collection",
			Status:  "500 Internal Server Error",
		}))
	}
	defer paymentDetailsCol.Database().Client().Disconnect(context.TODO())

	// Insert payment details into MongoDB
	_, err = paymentDetailsCol.InsertOne(ctx.Context(), paymentDetails)
	if err != nil {
		return ctx.JSON(commonutils.CreateFailureResponse(&common_responses.FailureResponse{
			Message: "Failed to save payment details",
			Status:  "500 Internal Server Error",
		}))
	}

	// Return the success response
	return ctx.JSON(commonutils.CreateSuccessResponse(&common_responses.SuccessResponse{
		Message: "Payment verified and details saved successfully",
		Status:  "200 OK",
		Data:    paymentDetails,
	}))
}

// verifySignature verifies the Razorpay payment signature
func verifySignature(verification *paymentModel.PaymentVerification) bool {
	apiSecret := commonutils.LoadEnv("RAZORPAY_API_SECRET")
	data := verification.OrderID + "|" + verification.PaymentID
	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write([]byte(data))
	expectedSignature := hex.EncodeToString(h.Sum(nil))
	return expectedSignature == verification.Signature
}
