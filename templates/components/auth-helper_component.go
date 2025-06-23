package components

import (
	"context"

	"github.com/nelsonmarro/bookings/config"
)

func IsUserAuthenticated(ctx context.Context) bool {
	app := config.GetConfigInstance()

	// Attempt to retrieve "user_id" from the session, expecting an int.
	sessionValue := app.Session.Get(ctx, "user_id")

	userID, ok := sessionValue.(int)
	if !ok {
		return false
	}

	isAuthenticated := userID > 0

	return isAuthenticated
}
