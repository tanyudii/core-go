package errutil

var (
	ErrAuthUnauthenticated      = NewUnauthenticatedError("unauthenticated")
	ErrAuthPermissionNotAllowed = NewUnauthorizedError("permission is not allowed")
	ErrAuthScopeNotAllowed      = NewUnauthorizedError("scope is not allowed")
	ErrAuthUserTypeNotAllowed   = NewUnauthorizedError("user type is not allowed")
)
