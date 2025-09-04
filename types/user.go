package types

type CreatePinBody struct {
	Pin             string `json:"pin" binding:"required"`
	PinConfirmation string `json:"pin_confirmation" binding:"required"`
}
