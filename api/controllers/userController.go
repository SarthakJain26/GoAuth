package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/SarthakJain26/GoAuth/api/models"
	"github.com/SarthakJain26/GoAuth/api/responses"
	"github.com/SarthakJain26/GoAuth/utils"
)

// UserSignUp controller for creating new users
func (a *App) UserSignUp(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Registered successfully"}
	user := &models.User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	usr, _ := user.GetUser(a.DB)
	if usr != nil {
		resp["status"] = "failed"
		resp["message"] = "User already registered, please login"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	user.Prepare() //To strip the white spaces

	err = user.Validate("") // default where all fields(email, lastname, firstname, password) are validated
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	userCreated, err := user.SaveUser(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["user"] = userCreated
	responses.JSON(w, http.StatusCreated, resp)
	return
}

// Login signs in users
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "logged in"}

	user := &models.User{}
	body, err := ioutil.ReadAll(r.Body) // read user input from request
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user.Prepare() //To strip the white spaces

	err = user.Validate("login") // fields(email, password) are validated
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	usr, err := user.GetUser(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if usr == nil { // user is not registered
		resp["status"] = "failed"
		resp["message"] = "Login failed, please signup"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	err = models.CheckPasswordHash(user.Password, usr.Password)
	if err != nil {
		resp["status"] = "failed"
		resp["message"] = "Login failed, please try again"
		responses.JSON(w, http.StatusForbidden, resp)
		return
	}
	token, err := utils.EncodeAuthToken(usr.ID)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["token"] = token
	responses.JSON(w, http.StatusOK, resp)
	return
}

// UpdateUserController updates users details
func (a *App) UpdateUserController(w http.ResponseWriter, r *http.Request) {
	// Default response
	var resp = map[string]interface{}{"status": "success", "message": "Details updated Successfully"}

	// Read request from user and check if the request is valid
	body, err := ioutil.ReadAll(r.Body)
	// If not vaild return suitable response
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user := &models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// If valid, update the user details and return a suitable response
	user, err = user.UpdateUser(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["user"] = user
	responses.JSON(w, http.StatusOK, resp)
	return
}

// DeactivateUserController deactivates users account
func (a *App) DeleteOrDeactivateUserController(w http.ResponseWriter, r *http.Request) {
	// Default response
	var resp = map[string]interface{}{"status": "success", "message": "User deactivated successfully"}

	// Read request from user and check if the request is valid
	body, err := ioutil.ReadAll(r.Body) // read user input from request
	// If not vaild return suitable response
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user := &models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Fetching route from URL
	reqURL := r.URL.String()
	isDelete := reqURL == "/delete"
	if isDelete {
		resp["message"] = "User deleted successfully"
	}

	// If valid, update the user details and return a suitable response
	err = user.DeleteOrDeactivateUser(a.DB, isDelete)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, resp)
	return

}
