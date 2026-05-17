package errs

import "errors"

var ErrInvalidCredential = errors.New(
	"invalid credential",
)

var ErrEmailNotFound = errors.New(
	"email not found",
)
var ErrExistingEmail = errors.New(
	"email has been registered",
)

var ErrInternalServer = errors.New(
	"internal server error",
)
