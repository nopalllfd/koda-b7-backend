package errs

import "errors"

var ErrInternalServer = errors.New(
	"internal server error",
)

var (
	ErrProfileNotFound  = errors.New("profile not found")
	ErrPhoneAlreadyUsed = errors.New("phone number already registered")
	ErrInvalidInput     = errors.New("invalid input")

	ErrInvalidCredential = errors.New("invalid email or password")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrEmailNotFound     = errors.New("email not found")
	ErrExistingEmail     = errors.New("email already registered")

	ErrUserAlreadyHasWallet = errors.New("user already has wallet")
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidBalance       = errors.New("invalid balance")
	ErrTimeoutOrCanceled    = errors.New("request timed out or canceled")

	ErrPINNotSet  = errors.New("pin is not set")
	ErrInvalidPin = errors.New("invalid pin")

	ErrWalletNotFound      = errors.New("wallet not found")
	ErrTransactionFailed   = errors.New("transaction failed")
	ErrTopupFailed         = errors.New("topup failed")
	ErrUpdateBalanceFailed = errors.New("failed to update balance")

	ErrSameWalletTransfer  = errors.New("cannot transfer to the same wallet")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrTransferFailed      = errors.New("transfer failed")
)
