package loginModel

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterReq struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterNewUserOrResetPasswordReq struct {
	Email      string `json:"email"`
	Rechaptcha string `json:"rechaptcha"`
}

type ForgotPasswordOTPRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}
type RegisterOTPRequest struct {
	Email    string `json:"email"`
	OTP      string `json:"otp"`
	Password string `json:"password"`
}

type ResetPasswordRequest struct {
	Email              string `json:"email"`
	NewPassword        string `json:"newPassword"`
	ConfirmNewPassword string `json:"confirmNewPassword"`
}

type SetPasswordReq struct {
	AuthToken string `json:"authToken"`
	Password  string `json:"assword"`
}
