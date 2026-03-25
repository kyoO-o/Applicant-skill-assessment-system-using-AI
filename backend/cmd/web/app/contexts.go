package app

type contextKey string

const (
	ContextKeyIsAuth         = contextKey("isAuthenticated")
	ContextKeyAuthCustomer   = contextKey("authenticatedCustomer")
	ContextKeyChosenCustomer = contextKey("chosenCustomer")
	ContextKeyChosenDocument = contextKey("chosenDocument")
	ContextKeyChosenPage     = contextKey("chosenPage")
)
