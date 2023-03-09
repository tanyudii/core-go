package errutil

var (
	ErrAuthUnauthenticated      = NewUnauthenticatedError("unauthenticated")
	ErrAuthPermissionNotAllowed = NewUnauthorizedError("permission is not allowed")
	ErrAuthScopeNotAllowed      = NewUnauthorizedError("scope is not allowed")
	ErrAuthScopeNotConfigured   = NewUnauthorizedError("user scope is not configured")
)
