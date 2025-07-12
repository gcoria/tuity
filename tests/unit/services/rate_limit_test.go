package services

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"tuity/internal/adapters/driving/http/middleware"
	"tuity/tests/testutils"

	"github.com/gin-gonic/gin"
)

// Test constants
const (
	testUserID1 = "user-123"
	testUserID2 = "user-456"
)

func Test_TokenBucket_Allow(t *testing.T) {
	// Create a token bucket with capacity 2, refill rate 2 per minute
	bucket := middleware.NewTokenBucket(2, 2)

	if !bucket.Allow() {
		t.Error("First request should be allowed")
	}

	if !bucket.Allow() {
		t.Error("Second request should be allowed")
	}

	if bucket.Allow() {
		t.Error("Third request should be denied (bucket empty)")
	}

	if bucket.Allow() {
		t.Error("Fourth request should be denied (bucket still empty)")
	}
}

func Test_TokenBucket_Refill(t *testing.T) {
	// Create a token bucket with capacity 2, refill rate 60 per minute (1 per second)
	bucket := middleware.NewTokenBucket(2, 60)

	bucket.Allow()
	bucket.Allow()

	if bucket.Allow() {
		t.Error("Request should be denied (bucket empty)")
	}

	time.Sleep(1 * time.Second)

	if !bucket.Allow() {
		t.Error("Request should be allowed after refill")
	}
}

func Test_TokenBucket_MaxCapacity(t *testing.T) {
	// Create a token bucket with capacity 2, refill rate 120 per minute (2 per second)
	bucket := middleware.NewTokenBucket(2, 120)

	bucket.Allow()

	time.Sleep(100 * time.Millisecond)

	if !bucket.Allow() {
		t.Error("Should allow first request (have ~1.2 tokens)")
	}
	if bucket.Allow() {
		t.Error("Should deny second request (only ~0.2 tokens left)")
	}

	time.Sleep(2 * time.Second)

	if !bucket.Allow() {
		t.Error("Should allow first request after full refill")
	}
	if !bucket.Allow() {
		t.Error("Should allow second request after full refill")
	}
	if bucket.Allow() {
		t.Error("Should deny third request (capacity is 2)")
	}
}

func Test_RateLimit_NoUserID(t *testing.T) {
	router := testutils.SetupTestRouter(middleware.TweetCreateRateLimit(5))

	response := testutils.MakeRequest(router, "")

	testutils.AssertStatusCode(t, response, http.StatusOK, "Request without user ID should be allowed")
}

func Test_RateLimit_WithinLimit(t *testing.T) {
	router := testutils.SetupTestRouter(
		middleware.TweetCreateRateLimit(3),
		middleware.ErrorHandler(),
	)

	responses := testutils.MakeMultipleRequests(router, testUserID1, 3)

	for i, response := range responses {
		testutils.AssertStatusCode(t, response, http.StatusOK, fmt.Sprintf("Request %d should be allowed", i+1))
	}
}

func Test_RateLimit_ExceedsLimit(t *testing.T) {
	router := testutils.SetupTestRouter(
		middleware.TweetCreateRateLimit(2),
		middleware.ErrorHandler(),
	)

	// Make 2 requests (within limit)
	responses := testutils.MakeMultipleRequests(router, testUserID1, 2)

	for i, response := range responses {
		testutils.AssertStatusCode(t, response, http.StatusOK, fmt.Sprintf("Request %d should be allowed", i+1))
	}

	// Make 3rd request (should be denied)
	response := testutils.MakeRequest(router, testUserID1)
	testutils.AssertStatusCode(t, response, http.StatusTooManyRequests, "Request 3 should be denied")
}

func Test_RateLimit_DifferentUsers(t *testing.T) {
	router := testutils.SetupTestRouter(
		middleware.TweetCreateRateLimit(2),
		middleware.ErrorHandler(),
	)

	user1Responses := testutils.MakeMultipleRequests(router, testUserID1, 2)
	user2Responses := testutils.MakeMultipleRequests(router, testUserID2, 2)

	for i, response := range user1Responses {
		testutils.AssertStatusCode(t, response, http.StatusOK, fmt.Sprintf("User 1 request %d should be allowed", i+1))
	}

	for i, response := range user2Responses {
		testutils.AssertStatusCode(t, response, http.StatusOK, fmt.Sprintf("User 2 request %d should be allowed", i+1))
	}

	user1Response := testutils.MakeRequest(router, testUserID1)
	user2Response := testutils.MakeRequest(router, testUserID2)

	testutils.AssertStatusCode(t, user1Response, http.StatusTooManyRequests, "User 1 request 3 should be denied")
	testutils.AssertStatusCode(t, user2Response, http.StatusTooManyRequests, "User 2 request 3 should be denied")
}

func Test_RateLimit_DifferentOperations(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tweetRateLimit := middleware.TweetCreateRateLimit(2)
	followRateLimit := middleware.FollowOperationRateLimit(2)

	router := gin.New()
	router.Use(middleware.ErrorHandler())

	router.GET("/tweets", tweetRateLimit, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "tweet success"})
	})

	router.GET("/follow", followRateLimit, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "follow success"})
	})

	tweetResponses := make([]*httptest.ResponseRecorder, 2)
	for i := 0; i < 2; i++ {
		tweetResponses[i] = testutils.MakeRequestToPath(router, "/tweets", testUserID1)
		testutils.AssertStatusCode(t, tweetResponses[i], http.StatusOK, fmt.Sprintf("Tweet request %d should be allowed", i+1))
	}

	followResponses := make([]*httptest.ResponseRecorder, 2)
	for i := 0; i < 2; i++ {
		followResponses[i] = testutils.MakeRequestToPath(router, "/follow", testUserID1)
		testutils.AssertStatusCode(t, followResponses[i], http.StatusOK, fmt.Sprintf("Follow request %d should be allowed", i+1))
	}

	tweetResponse := testutils.MakeRequestToPath(router, "/tweets", testUserID1)
	followResponse := testutils.MakeRequestToPath(router, "/follow", testUserID1)

	testutils.AssertStatusCode(t, tweetResponse, http.StatusTooManyRequests, "Tweet request 3 should be denied")
	testutils.AssertStatusCode(t, followResponse, http.StatusTooManyRequests, "Follow request 3 should be denied")
}
