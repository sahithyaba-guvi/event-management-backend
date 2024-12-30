package login

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"

	mongoSetup "em_backend/configs/mongo"
	redisSetup "em_backend/configs/redis"
	common_responses "em_backend/responses/common"
	loginRespo "em_backend/responses/login"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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
