package kafekoding_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aZ4ziL/kafekoding_api/models"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// A handler for handling new user registration requests.
func signUpHandler(w http.ResponseWriter, r *http.Request) {
	payloads := struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
		Username  string `json:"username" validate:"required"`
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required"`
	}{}

	// set content-type to json
	w.Header().Set("Content-Type", "application/json")

	// parse request body to json payloads;
	err := json.NewDecoder(r.Body).Decode(&payloads)
	if err != nil {
		// if request body is not json.
		data := map[string]interface{}{
			"status":  "error",
			"message": "Payload yang anda gunakan tidak di ijinkan.",
		}
		responseJSON(w, http.StatusBadRequest, data)
		return
	}

	// check validate for field
	validate = validator.New()
	err = validate.Struct(&payloads)
	if err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			if err.ActualTag() == "required" {
				// TODO: translate field name for indonesia language.
				errorMessages = append(errorMessages, fmt.Sprintf("Bidang dengan nama `%s` ini dibutuhkan.", err.Field()))
			} else if err.ActualTag() == "email" {
				errorMessages = append(errorMessages, fmt.Sprintf("Bidang dengan nama `%s` ini harus berisi email yang valid", err.Field()))
			}
		}
		data := map[string]interface{}{
			"status":  "error",
			"message": "Mohon perbaiki sesalahan di bawah ini.",
			"errors":  errorMessages,
		}
		responseJSON(w, http.StatusBadRequest, data)
		return
	}

	// save data payload to models.User
	user := models.User{
		FirstName: payloads.FirstName,
		LastName:  payloads.LastName,
		Username:  payloads.Username,
		Email:     payloads.Email,
		Password:  payloads.Password,
	}
	// save to db
	err = models.CreateNewUser(&user)
	if err != nil {
		if strings.Contains(err.Error(), "users_username_key") {
			responseJSON(w, http.StatusBadRequest, map[string]interface{}{
				"status":  "error",
				"message": "Username yang anda gunakan telah terdaftar di akun lain, silahkan gunakan username yang berbeda.",
			})
			return
		}
		if strings.Contains(err.Error(), "users_email_key") {
			responseJSON(w, http.StatusBadRequest, map[string]interface{}{
				"status":  "error",
				"message": "Email yang anda gunakan telah terdaftar di akun lain, silahkan gunakan email yang berbeda.",
			})
			return
		}
		responseJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// success reponse
	responseJSON(w, http.StatusCreated, user)
}

// getTokenHandler is handler for handling request for get token.
func getTokenHandler(w http.ResponseWriter, r *http.Request) {
	payloads := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	w.Header().Set("Content-Type", "application/json")

	// decode the request.Body to payloads
	err := json.NewDecoder(r.Body).Decode(&payloads)
	if err != nil {
		responseJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Payload yang anda gunakan tidak diijinkan.",
		})
		return
	}

	var user models.User
	if strings.Contains(payloads.Username, "@") {
		user, err = models.GetUserByEmail(payloads.Username)
		if err != nil {
			responseJSON(w, http.StatusBadRequest, map[string]interface{}{
				"status":  "error",
				"message": "Username, email atau katasandi yang anda masukkan salah.",
			})
			return
		}
	} else {
		user, err = models.GetUserByUsername(payloads.Username)
		if err != nil {
			responseJSON(w, http.StatusBadRequest, map[string]interface{}{
				"status":  "error",
				"message": "Username, email atau katasandi yang anda masukkan salah.",
			})
			return
		}
	}

	// check password
	if !models.DecryptionPassword(user.Password, payloads.Password) {
		responseJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Username, email atau katasandi yang anda masukkan salah.",
		})
		return
	}

	credential := credential{
		ID: user.ID,
	}

	token, err := generateNewToken(credential)
	if err != nil {
		responseJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// success response
	responseJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Harap jangan bagikan kode akses token ini kepada siapapun.",
		"token":   token,
	})
}
