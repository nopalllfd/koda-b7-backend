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

	ErrUserAlreadyHasWallet = errors.New("user sudah memiliki dompet")
	ErrUserNotFound         = errors.New("user tidak ditemukan")
	ErrInvalidBalance       = errors.New("saldo tidak valid")
	ErrTimeoutOrCanceled    = errors.New("proses dibatalkan atau waktu habis")

	ErrPINNotSet  = errors.New("user belum mengatur PIN")
	ErrInvalidPin = errors.New("pin salah")
)
