package app

type contextKey string

const (
	ContextKeyIsAuth         = contextKey("isAuthenticated")
	ContextKeyAuthCustomer   = contextKey("authenticatedCustomer")
	ContextKeyChosenCustomer = contextKey("chosenCustomer")
	ContextKeyChosenCompany  = contextKey("chosenCompany")
	ContextKeyChosenJob      = contextKey("chosenJob")
	ContextKeyChosenDocument = contextKey("chosenDocument")
	ContextKeyChosenPage     = contextKey("chosenPage")
)
