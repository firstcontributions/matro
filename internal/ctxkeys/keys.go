package ctxkeys

// ContextKey is the type user for setting keys in context.WithVal
type ContextKey int

const (
	// Spinner used as a key while adding spinner.Spinner to ctx
	Spinner ContextKey = iota
)
