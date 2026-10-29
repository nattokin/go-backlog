package backlog_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	backlog "github.com/nattokin/go-backlog"
)

// Example_errorHandling demonstrates telling the returned errors apart with
// errors.As. Check *OAuthError before *APIResponseError, because an
// *OAuthError wraps an *APIResponseError.
func Example_errorHandling() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
		backlog.WithDoer(doerNotFound),
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	describe := func(err error) {
		var (
			oauthErr      *backlog.OAuthError
			apiErr        *backlog.APIResponseError
			validationErr *backlog.ValidationError
			optionKeyErr  *backlog.InvalidOptionKeyError
		)
		switch {
		case errors.As(err, &oauthErr):
			fmt.Println("OAuth error:", oauthErr.Code())
		case errors.As(err, &apiErr):
			fmt.Println("API error: status", apiErr.StatusCode())
			for _, e := range apiErr.Errors() {
				fmt.Println("message:", e.Message)
			}
		case errors.As(err, &validationErr):
			fmt.Printf("validation error: %s: %s\n", validationErr.Target(), validationErr.Message())
		case errors.As(err, &optionKeyErr):
			fmt.Println("invalid option:", optionKeyErr.InvalidKey())
		case err != nil:
			fmt.Println("error:", err)
		}
	}

	ctx := context.Background()

	// The server responds with 404 Not Found.
	_, err = c.Wiki.One(ctx, 1)
	describe(err)

	// The argument is invalid, so no request is sent.
	_, err = c.Wiki.One(ctx, 0)
	describe(err)

	// The option is not accepted by this method, so no request is sent.
	_, err = c.Wiki.List(ctx, "MYPROJECT", c.Issue.Option.WithCount(10))
	describe(err)
	// Output:
	// API error: status 404
	// message: No such resource.
	// validation error: wikiID: invalid wikiID: must not be less than 1
	// invalid option: count
}

// ExampleAPIResponseError_RateLimitStatus demonstrates reading the rate limit
// status after a 429 Too Many Requests response.
func ExampleAPIResponseError_RateLimitStatus() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
		backlog.WithDoer(doerRateLimited),
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	_, err = c.Wiki.List(context.Background(), "MYPROJECT")

	var apiErr *backlog.APIResponseError
	if errors.As(err, &apiErr) {
		if rl := apiErr.RateLimitStatus(); rl != nil {
			fmt.Println("status:", apiErr.StatusCode())
			fmt.Println("limit:", rl.Limit)
			fmt.Println("remaining:", rl.Remaining)
			fmt.Println("reset:", rl.Reset.Format(time.RFC3339))
			// Wait until rl.Reset before sending the next request.
		}
	}
	// Output:
	// status: 429
	// limit: 600
	// remaining: 0
	// reset: 2025-01-01T00:00:00Z
}
