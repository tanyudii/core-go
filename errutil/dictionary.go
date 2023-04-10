package errutil

var (
	ErrAuthUnauthenticated      = NewUnauthenticatedErrorWithName("unauthenticated", "UNAUTHENTICATED")
	ErrAuthPermissionNotAllowed = NewUnauthorizedErrorWithName("permission is not allowed", "PERMISSION_NOT_ALLOWED")
	ErrAuthScopeNotAllowed      = NewUnauthorizedErrorWithName("scope is not allowed", "SCOPE_NOT_ALLOWED")
	ErrAuthUserTypeNotAllowed   = NewUnauthorizedErrorWithName("user type is not allowed", "USER_TYPE_NOT_ALLOWED")
)
