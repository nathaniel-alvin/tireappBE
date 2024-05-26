package tireappbe

import (
	"context"

	"github.com/nathaniel-alvin/tireappBE/types"
)

type contextKey int

// different with ecom, might cause a problem
const userContextKey = contextKey(iota + 1)

// return a new context with a given user
func NewContextWithUser(ctx context.Context, user *types.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// return the current logged in user
func UserFromContext(ctx context.Context) *types.User {
	user, _ := ctx.Value(userContextKey).(*types.User)
	return user
}

// return ID of the current logged in user. return -1 if not logged in
func UserIDFromContext(ctx context.Context) int {
	if user := UserFromContext(ctx); user != nil {
		return user.ID
	}
	return -1
}
