# Test Utilities

This package contains shared utilities and mocks for testing across the Tuity project.

## Mocks

The `mocks` package provides mock implementations of repository interfaces for testing:

### Available Mocks

- **MockUserRepository**: Mock implementation of `UserRepository` interface
- **MockTweetRepository**: Mock implementation of `TweetRepository` interface

### Usage

```go
package services

import (
    "testing"
    "tuity/tests/testutils/mocks"
)

func TestMyService(t *testing.T) {
    // Create mock repositories
    userRepo := mocks.NewMockUserRepository()
    tweetRepo := mocks.NewMockTweetRepository()

    // Configure mock behavior
    userRepo.ShouldFailSave = true  // Simulate repository failure

    // Use mocks in your tests
    service := services.NewMyService(userRepo, tweetRepo)
    // ... test logic
}
```

### Mock Configuration

Mocks can be configured to simulate different scenarios:

- **ShouldFailSave**: Set to `true` to make the `Save` method return an error

## Test Helpers

The `testutils` package provides common assertion helpers:

### Available Helpers

- **AssertError**: Checks if an error is of the expected type
- **AssertNoError**: Checks that no error occurred
- **AssertValidTweet**: Validates a tweet has expected properties
- **RepeatChar**: Creates a string by repeating a character n times

### Usage

```go
package services

import (
    "testing"
    "tuity/tests/testutils"
)

func TestMyFunction(t *testing.T) {
    result, err := myFunction()

    testutils.AssertNoError(t, err)
    testutils.AssertValidTweet(t, result, "user-123", "Hello world!")
}
```

## HTTP Test Utilities

The `testutils` package also provides HTTP testing utilities for testing HTTP handlers and middleware:

### Available HTTP Utilities

- **SetupTestRouter**: Creates a test router with optional middleware
- **MakeRequest**: Makes HTTP GET requests with optional user ID header
- **MakeRequestToPath**: Makes HTTP GET requests to specific paths
- **AssertStatusCode**: Asserts HTTP response status codes
- **MakeMultipleRequests**: Makes multiple HTTP requests

### Usage

```go
package services

import (
    "testing"
    "net/http"
    "tuity/tests/testutils"
    "tuity/internal/adapters/driving/http/middleware"
)

func TestMyMiddleware(t *testing.T) {
    // Setup test router with middleware
    router := testutils.SetupTestRouter(
        middleware.MyMiddleware(),
        middleware.ErrorHandler(),
    )

    // Make requests
    response := testutils.MakeRequest(router, "user-123")

    // Assert response
    testutils.AssertStatusCode(t, response, http.StatusOK, "Request should succeed")
}
```

## Benefits

- **Reusable**: Eliminates code duplication across test files
- **Consistent**: Provides consistent mock behavior and assertions
- **Maintainable**: Changes to mock behavior only need to be made in one place
- **Extensible**: Easy to add new mocks and helpers as needed
