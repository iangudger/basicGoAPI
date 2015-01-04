package api

const (
	duplicateParameters    = "Request contains duplicate parameters."
	invalidLoginParameters = "Invalid login parameters."
	missingAction          = "Action is missing from request."
	invalidAction          = "Invalid action."
	invalidEmail           = "The provided email address is not valid."
	registrationSuccess    = "Registration successful."
	invalidSessionID       = "Invalid sessionid."
	unknownUser            = "There is no record of this user."
	passwordResetSuccess   = "Your password has been successfully reset. Please check your email for your new password."
	invalidPassword        = "Invalid password."
	passwordChangeSuccess  = "Password has been successfully changed."
	logoutSuccess          = "Logout successful."
)

type ApiError struct {
	Message string `json:"errorMessage"`
}

type ApiResult struct {
	Message string `json:"message"`
}

type ApiData struct {
	Data interface{} `json:"data"`
}

type Session struct {
	ID string `json:"sessionid"`
}

type Username struct {
	Username string `json:"username"`
}
