package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"tuity/internal/adapters/driving/http/dto"
	"tuity/internal/app"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	container := app.NewContainer()
	return app.SetupRouter(container)
}

func TestHealthCheck(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", response["status"])
	}
}

func Test_User_PostGet(t *testing.T) {
	router := setupTestRouter()

	// Create a user
	createUserReq := dto.CreateUserRequest{
		Username:    "diego",
		DisplayName: "Diego Maradona",
	}

	reqBody, _ := json.Marshal(createUserReq)
	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var createdUser dto.UserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &createdUser); err != nil {
		t.Errorf("Failed to unmarshal user response: %v", err)
	}

	if createdUser.Username != "diego" {
		t.Errorf("Expected username 'diego', got %s", createdUser.Username)
	}
	if createdUser.DisplayName != "Diego Maradona" {
		t.Errorf("Expected display name 'Diego Maradona', got %s", createdUser.DisplayName)
	}
	if createdUser.ID == "" {
		t.Error("Expected user ID to be set")
	}

	req = httptest.NewRequest("GET", "/api/v1/users/"+createdUser.ID, nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var retrievedUser dto.UserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &retrievedUser); err != nil {
		t.Errorf("Failed to unmarshal user response: %v", err)
	}

	if retrievedUser.ID != createdUser.ID {
		t.Errorf("Expected user ID %s, got %s", createdUser.ID, retrievedUser.ID)
	}
}

func Test_Tweet_PostGet(t *testing.T) {
	router := setupTestRouter()

	//create a user
	createUserReq := dto.CreateUserRequest{
		Username:    "lio",
		DisplayName: "Lionel Messi",
	}

	reqBody, _ := json.Marshal(createUserReq)
	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var createdUser dto.UserResponse
	json.Unmarshal(w.Body.Bytes(), &createdUser)

	// Create a tweet
	createTweetReq := dto.CreateTweetRequest{
		Content: "Hello, Tuity world! 🌍",
	}

	reqBody, _ = json.Marshal(createTweetReq)
	req = httptest.NewRequest("POST", "/api/v1/tweets", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", createdUser.ID)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var createdTweet dto.TweetResponse
	if err := json.Unmarshal(w.Body.Bytes(), &createdTweet); err != nil {
		t.Errorf("Failed to unmarshal tweet response: %v", err)
	}

	if createdTweet.Content != "Hello, Tuity world! 🌍" {
		t.Errorf("Expected content 'Hello, Tuity world! 🌍', got %s", createdTweet.Content)
	}
	if createdTweet.UserID != createdUser.ID {
		t.Errorf("Expected user ID %s, got %s", createdUser.ID, createdTweet.UserID)
	}

	//get the tweet
	req = httptest.NewRequest("GET", "/api/v1/tweets/"+createdTweet.ID, nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var retrievedTweet dto.TweetResponse
	if err := json.Unmarshal(w.Body.Bytes(), &retrievedTweet); err != nil {
		t.Errorf("Failed to unmarshal tweet response: %v", err)
	}

	if retrievedTweet.ID != createdTweet.ID {
		t.Errorf("Expected tweet ID %s, got %s", createdTweet.ID, retrievedTweet.ID)
	}
}

func Test_Following(t *testing.T) {
	router := setupTestRouter()

	//create two users
	users := []dto.UserResponse{}
	usernames := []string{"diego", "lio"}
	displayNames := []string{"Diego Maradona", "Lionel Messi"}

	for i, username := range usernames {
		createUserReq := dto.CreateUserRequest{
			Username:    username,
			DisplayName: displayNames[i],
		}

		reqBody, _ := json.Marshal(createUserReq)
		req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var user dto.UserResponse
		json.Unmarshal(w.Body.Bytes(), &user)
		users = append(users, user)
	}

	diego := users[0]
	lio := users[1]

	//diego follows lio
	req := httptest.NewRequest("POST", "/api/v1/users/"+lio.ID+"/follow", nil)
	req.Header.Set("X-User-ID", diego.ID)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	//check diego's following list
	req = httptest.NewRequest("GET", "/api/v1/users/"+diego.ID+"/following", nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var diegoFollowing []dto.FollowResponse
	if err := json.Unmarshal(w.Body.Bytes(), &diegoFollowing); err != nil {
		t.Errorf("Failed to unmarshal following response: %v", err)
	}

	if len(diegoFollowing) != 1 {
		t.Errorf("Expected 1 following, got %d", len(diegoFollowing))
	}
	if diegoFollowing[0].FollowedID != lio.ID {
		t.Errorf("Expected diego to follow lio (%s), got %s", lio.ID, diegoFollowing[0].FollowedID)
	}

	//check lio's followers list
	req = httptest.NewRequest("GET", "/api/v1/users/"+lio.ID+"/followers", nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var lioFollowers []dto.FollowResponse
	if err := json.Unmarshal(w.Body.Bytes(), &lioFollowers); err != nil {
		t.Errorf("Failed to unmarshal followers response: %v", err)
	}

	if len(lioFollowers) != 1 {
		t.Errorf("Expected 1 follower, got %d", len(lioFollowers))
	}
	if lioFollowers[0].FollowerID != diego.ID {
		t.Errorf("Expected lio to be followed by el diego (%s), got %s", diego.ID, lioFollowers[0].FollowerID)
	}
}

func Test_Timeline_Generation(t *testing.T) {
	router := setupTestRouter()

	// Create two users
	users := []dto.UserResponse{}
	usernames := []string{"diego", "lio"}
	displayNames := []string{"Diego Maradona", "Lionel Messi"}

	for i, username := range usernames {
		createUserReq := dto.CreateUserRequest{
			Username:    username,
			DisplayName: displayNames[i],
		}

		reqBody, _ := json.Marshal(createUserReq)
		req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var user dto.UserResponse
		json.Unmarshal(w.Body.Bytes(), &user)
		users = append(users, user)
	}

	diego := users[0]
	lio := users[1]

	//diego follows lio
	req := httptest.NewRequest("POST", "/api/v1/users/"+lio.ID+"/follow", nil)
	req.Header.Set("X-User-ID", diego.ID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	//lio creates a tweet
	createTweetReq := dto.CreateTweetRequest{
		Content: "Lio's first tweet!",
	}

	reqBody, _ := json.Marshal(createTweetReq)
	req = httptest.NewRequest("POST", "/api/v1/tweets", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", lio.ID)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	//diego creates a tweet
	createTweetReq = dto.CreateTweetRequest{
		Content: "Diego's first tweet!",
	}

	reqBody, _ = json.Marshal(createTweetReq)
	req = httptest.NewRequest("POST", "/api/v1/tweets", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", diego.ID)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	//get diego's timeline
	req = httptest.NewRequest("GET", "/api/v1/users/"+diego.ID+"/timeline", nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var timeline dto.TimelineResponse
	if err := json.Unmarshal(w.Body.Bytes(), &timeline); err != nil {
		t.Errorf("Failed to unmarshal timeline response: %v", err)
	}

	//diego's timeline should contain both his tweet and lio's tweet
	if len(timeline.Tweets) != 2 {
		t.Errorf("Expected 2 tweets in timeline, got %d", len(timeline.Tweets))
	}

	// Check that both tweets are present
	foundDiegoTweet := false
	foundLioTweet := false
	for _, tweet := range timeline.Tweets {
		if tweet.Content == "Diego's first tweet!" {
			foundDiegoTweet = true
		}
		if tweet.Content == "Lio's first tweet!" {
			foundLioTweet = true
		}
	}

	if !foundDiegoTweet {
		t.Error("Diego's tweet should be in his timeline")
	}
	if !foundLioTweet {
		t.Error("Lio's tweet should be in diego's timeline (he follows him)")
	}
}

func Test_Tweet_CharacterLimit(t *testing.T) {
	router := setupTestRouter()

	// Create a user
	createUserReq := dto.CreateUserRequest{
		Username:    "enzo",
		DisplayName: "Enzo Perez",
	}

	reqBody, _ := json.Marshal(createUserReq)
	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var createdUser dto.UserResponse
	json.Unmarshal(w.Body.Bytes(), &createdUser)

	// Try to create a tweet that exceeds 280 characters
	longContent := ""
	for i := 0; i < 281; i++ {
		longContent += "a"
	}

	createTweetReq := dto.CreateTweetRequest{
		Content: longContent,
	}

	reqBody, _ = json.Marshal(createTweetReq)
	req = httptest.NewRequest("POST", "/api/v1/tweets", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", createdUser.ID)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for tweet exceeding character limit, got %d", w.Code)
	}

	var errorResponse dto.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
		t.Errorf("Failed to unmarshal error response: %v", err)
	}

	if errorResponse.Error != "validation_error" {
		t.Errorf("Expected validation error, got %s", errorResponse.Error)
	}
}
