package login

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	credendials "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/go-redis/redis"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"

	mongoSetup "em_backend/configs/mongo"
	redisSetup "em_backend/configs/redis"
	commonutils "em_backend/library/common"
	common_responses "em_backend/responses/common"
	loginRespo "em_backend/responses/login"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
)

func CheckSession(auth string) loginRespo.LoginResp {
	var resp loginRespo.LoginResp
	r, err := redisSetup.ConnectToRedis()
	if err != nil {
		fmt.Println(err)
	}
	result := r.Get(auth)
	value, err := result.Result()
	if err != nil {
		fmt.Println(err)
	}
	var userInfo common_responses.LoginDetails
	err = json.Unmarshal([]byte(value), &userInfo)
	if err != nil {
		fmt.Println(err)
	}
	if userInfo.Hash != "" {
		resp.Access = true
		resp.UserInfo = userInfo
	} else {
		resp.Access = false
	}
	return resp

}

// Helper function to generate a JWT token (no expiry)
func GenerateAuthToken(userHash string) (string, error) {
	claims := jwt.MapClaims{
		"userHash": userHash,
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Signing the token with a secret key
	secretKey := "your_secret_key_here" // Use a strong secret key here
	authToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return authToken, nil
}

// Helper function to store auth token in Redis
func StoreAuthTokenInRedis(userRedisData common_responses.LoginDetails) error {
	// Get the Redis client by calling the ConnectToRedis function
	redisClient, err := redisSetup.ConnectToRedis() // Assuming this returns a Redis client
	if err != nil {
		return fmt.Errorf("error connecting to Redis: %v", err)
	}
	defer redisClient.Close()

	// Check if the userHash already exists in hashMapper
	oldAuthToken, err := redisClient.HGet("hashMapper", userRedisData.Hash).Result()
	if err == nil && oldAuthToken != "" {
		// If an old auth token exists, delete the associated auth:userInfo entry
		err = redisClient.Del(oldAuthToken).Err()
		if err != nil {
			return fmt.Errorf("error deleting old auth:userInfo entry: %v", err)
		}
	}

	// Serialize userRedisData to JSON
	userDataJSON, err := json.Marshal(userRedisData)
	if err != nil {
		return fmt.Errorf("error marshalling user data to JSON: %v", err)
	}

	// Store the authToken:userInfo mapping
	err = redisClient.Set(userRedisData.AuthToken, userDataJSON, 0).Err() // 0 means no expiration time
	if err != nil {
		return fmt.Errorf("error storing auth token in Redis: %v", err)
	}

	// Store the hash:auth mapping inside "hashMapper"
	err = redisClient.HSet("hashMapper", userRedisData.Hash, userRedisData.AuthToken).Err()
	if err != nil {
		return fmt.Errorf("error storing hash:auth mapping: %v", err)
	}

	return nil
}

func IsEmailRegistered(email string) (bool, error) {
	collectionName := "userData"
	filter := bson.M{"email": email}
	projection := bson.M{"_id": 1} // Only fetch the `_id` field for efficiency.

	result, err := mongoSetup.FindOneDoc(collectionName, filter, projection)
	if err != nil {
		fmt.Println(err)
		return false, fmt.Errorf("error querying database: %w", err)
	}

	// Check if a document was found.
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return false, nil // Email not found
		}
		return false, result.Err() // Some other error
	}

	return true, nil // Email exists
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func GenerateUserHash(email string) string {
	data := email + time.Now().String()
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func GenerateOTP() string {
	return fmt.Sprintf("%04d", rand.Intn(10000))
}

func SendRegisterOrForgotPasswordEmail(email string, mailType string) (bool, error) {
	r, err := redisSetup.ConnectToRedis()
	if err != nil {
		fmt.Println("Err connecting to redis", err)
	}
	defer r.Close()

	generatedPassword := GenerateOTP()

	expiryTime := 5 * time.Minute
	email = "dukone.contact@gmail.com"
	key := fmt.Sprintf("emailverification:%s", email)

	err = r.Set(key, generatedPassword, expiryTime).Err()

	if err != nil {
		fmt.Println("Cannot set generated password in redis")
	}

	const (
		sender    = "dukone.contact@gmail.com"
		awsRegion = "us-east-1"
	)

	actionMessage := ""
	if mailType == "forgotpassword" {
		actionMessage = "Please use the verification code below to reset your password:"
	} else {
		actionMessage = "Please use the verification code below to complete your registration on STUNI"
	}
	bodyHTML := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			body { font-family: Arial, sans-serif; color: #333; margin: 20px; }
			.container { max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd; border-radius: 10px; background-color: #f9f9f9; }
			h2 { color: #007BFF; }
			.code { font-weight: bold; font-size: 20px; color: #d9534f; }
			.footer { font-size: 12px; color: #777; margin-top: 20px; }
		</style>
	</head>
	<body>
		<div class="container">
			<h2>STUNI Email Verification</h2>
			<p>Hi,</p>
			<p>%s</p>
			<p class="code">%s</p>
			<p>This code will expire in <strong>%d minutes</strong>.</p>
			<p>If you did not request this code, please ignore this email.</p>
			<p>Best regards,</p>
			<p>The STUNI Team</p>
			<div class="footer">
				<p>Note: This is an automated email. Please do not reply.</p>
			</div>
		</div>
	</body>
	</html>`, actionMessage, string(generatedPassword), int(expiryTime/time.Minute))
	fmt.Println()
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credendials.NewStaticCredentials(commonutils.LoadEnv("AWS_KEY"), commonutils.LoadEnv("AWS_SKEY"), ""),
	})
	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
		return false, err
	}

	// Create SES service client
	svc := ses.New(sess)

	// Compose email input
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(email),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Data: aws.String(bodyHTML),
				},
			},
			Subject: &ses.Content{
				Data: aws.String("STUNI Email Verification"),
			},
		},
		Source: aws.String(sender),
	}

	// Send the email
	_, err = svc.SendEmail(input)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return false, err
	}
	if mailType == "forgotpassword" {
		log.Printf("Forgot password email sent successfully to %s", email)
	} else {
		log.Printf("Register account email sent successfully to %s", email)
	}
	return true, nil
}
func VerifyForgotPasswordOTP(email string, OTP string) (bool, error) {
	r, err := redisSetup.ConnectToRedis()
	if err != nil {
		fmt.Println(err)
	}
	defer r.Close()

	key := fmt.Sprintf("emailverification:%s", email)

	// Retrieve the stored OTP from Redis
	storedOTP, err := r.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			// Key does not exist, OTP expired or not set
			return false, fmt.Errorf("OTP has expired")
		}
		// Other Redis errors
		return false, fmt.Errorf("failed to retrieve OTP from Redis: %v", err)
	}

	if OTP != storedOTP {
		return false, fmt.Errorf("invalid OTP")
	}

	// If the OTP is valid, optionally delete it from Redis to prevent reuse
	err = r.Del(key).Err()
	if err != nil {
		return true, fmt.Errorf("failed to delete OTP from Redis: %v", err)
	}

	key = fmt.Sprintf("verificationDone:%s", email)
	// OTP is valid
	r.Set(key, "true", 0)
	return true, nil
}
