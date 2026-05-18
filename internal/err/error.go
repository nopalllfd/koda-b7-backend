package errs

import "errors"

var ErrInternalServer = errors.New(
	"internal server error",
)

var (
	ErrProfileNotFound  = errors.New("profil tidak ditemukan")
	ErrPhoneAlreadyUsed = errors.New("nomor telepon sudah terdaftar")
	ErrInvalidInput     = errors.New("data yang dimasukkan salah")

	ErrInvalidCredential = errors.New("invalid credential")
	ErrEmailNotFound     = errors.New("email not found")
	ErrExistingEmail     = errors.New("email has been registered")
)
