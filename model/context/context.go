package context

type ContextKey string

var (
	ContextKeyProperty    = ContextKey("property")
	ContextKeyPerson      = ContextKey("person")
	ContextKeyAccountType = ContextKey("accountType")
	ContextCustomClaims   = ContextKey("customClaims")
	ContextKeyUUID        = ContextKey("UUID")
	ContextKeyIP          = ContextKey("IP")
)
