package api

import (
	"log"
	"net/http"

	"github.com/iangudger/basicGoAPI/common"
	"github.com/iangudger/basicGoAPI/database"
)

var actions = map[string]func(http.ResponseWriter, *http.Request, database.DB) interface{}{}

func init() {
	actions["login"] = login
	actions["logout"] = logout
	actions["getUsername"] = getUsername
	actions["register"] = register
	actions["resetPassword"] = resetPassword
	actions["changePassword"] = changePassword
}

func login(res http.ResponseWriter, req *http.Request, dbconn database.DB) interface{} {
	params := req.URL.Query()

	// Make sure parameters are present
	if !checkParam(params, "username") ||
		!checkParam(params, "password") {
		common.InternalServerError(res, invalidLoginParameters)
		return nil
	}

	// Initiate login
	username := params["username"][0]
	password := params["password"][0]
	log.Printf("Initiating login for user: %s\n", username)
	sessionid, err := common.Login(dbconn, username, password)

	if err != nil {
		log.Printf("Login for user: %s failed with error message: %s\n", username, err.Error())
		return ApiError{err.Error()}
	} else {
		log.Printf("Login for user: %s succeeded. New sessionid: %s\n", username, sessionid)
		return Session{sessionid}
	}
}

func register(res http.ResponseWriter, req *http.Request, dbconn database.DB) interface{} {
	params := req.URL.Query()

	// Make sure parameters are present
	if !checkParam(params, "username") {
		common.InternalServerError(res, invalidEmail)
		return nil
	}

	username := params["username"][0]
	log.Printf("Initiating registration for user: %s\n", username)

	if !common.ValidEmail(username) {
		return ApiError{invalidEmail}
	}

	sessionid, err := common.Register(dbconn, username)
	if err != nil {
		return ApiError{err.Error()}
	}

	return Session{sessionid}
}

func getUsername(res http.ResponseWriter, req *http.Request, dbconn database.DB) interface{} {
	params := req.URL.Query()

	// Make sure parameters are present
	if !checkParam(params, "sessionid") {
		common.InternalServerError(res, invalidSessionID)
		return nil
	}

	sessionid := params["sessionid"][0]
	log.Printf("Looking up email with sessionid: %s\n", sessionid)

	username, err := common.GetEmail(dbconn, sessionid)
	if err != nil {
		log.Printf("Error retrieving email from sessionid (%s): %s\n", sessionid, err.Error())
		return ApiError{invalidSessionID}
	} else {
		log.Printf("Sessionid: %s is associated with the email: %s\n", sessionid, username)
		return Username{username}
	}
}

func resetPassword(res http.ResponseWriter, req *http.Request, dbconn database.DB) interface{} {
	params := req.URL.Query()

	// Make sure parameters are present
	if !checkParam(params, "username") {
		common.InternalServerError(res, invalidEmail)
		return nil
	}

	username := params["username"][0]
	log.Printf("Resetting password for user: %s\n", username)

	err := common.ResetPassword(dbconn, username)
	if err != nil {
		log.Printf("Error resetting password for email (%s): %s\n", username, err.Error())
		return ApiError{unknownUser}
	}

	return ApiResult{passwordResetSuccess}
}

func changePassword(res http.ResponseWriter, req *http.Request, dbconn database.DB) interface{} {
	params := req.URL.Query()

	// Make sure parameters are present
	if !checkParam(params, "sessionid") {
		common.InternalServerError(res, invalidSessionID)
		return nil
	}
	if !checkParam(params, "oldPassword") ||
		!checkParam(params, "newPassword") {
		common.InternalServerError(res, invalidPassword)
		return nil
	}

	sessionid := params["sessionid"][0]
	oldPassword := params["oldPassword"][0]
	newPassword := params["newPassword"][0]

	err := common.ChangePassword(dbconn, sessionid, oldPassword, newPassword)
	if err != nil {
		return ApiError{err.Error()}
	}

	return ApiResult{passwordChangeSuccess}
}

func logout(res http.ResponseWriter, req *http.Request, dbconn database.DB) interface{} {
	params := req.URL.Query()

	// Make sure parameters are present
	if !checkParam(params, "sessionid") {
		common.InternalServerError(res, invalidSessionID)
		return nil
	}

	sessionid := params["sessionid"][0]
	log.Printf("Logging out sessionid: %s\n", sessionid)

	err := common.Logout(dbconn, sessionid)
	if err != nil {
		log.Printf("Error logging out sessionid (%s): %s\n", sessionid, err.Error())
		return ApiError{invalidSessionID}
	}

	return ApiResult{logoutSuccess}
}
