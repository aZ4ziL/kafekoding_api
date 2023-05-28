package kafekoding_api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aZ4ziL/kafekoding_api/models"
	"github.com/go-playground/validator/v10"
)

// courseHandlerGET is handler for handling the course data with request type is GET.
func courseHandlerGET(w http.ResponseWriter, r *http.Request) {
	// if query id
	id := r.URL.Query().Get("id")
	if id != "" {
		courseHandlerDetailGET(w, r)
		return
	}

	courses := models.GetAllCourse(true)
	responseJSON(w, http.StatusOK, courses)
}

func courseHandlerDetailGET(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	idInt, _ := strconv.Atoi(id)

	course, err := models.GetCourseByID(idInt, true)
	if err != nil {
		responseJSON(w, http.StatusNotFound, map[string]interface{}{
			"status":  "error",
			"message": "Kursus yang anda tuju tidak dapat kami temukan.",
		})
		return
	}

	responseJSON(w, http.StatusOK, course)
}

// courseHandlerPOST is handle for create new course.
func courseHandlerPOST(w http.ResponseWriter, r *http.Request) {
	// check use is admin or not.
	// If not admin request stop.
	userCtx := r.Context().Value(&userAuth{}).(claims)
	user, _ := models.GetUserByID(userCtx.credential.ID)
	if !user.IsAdmin {
		responseJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Anda tidak memiliki akses untuk ini.",
		})
		return
	}

	// parse form data
	if err := r.ParseMultipartForm(1024); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	title := r.FormValue("title")
	file, header, err := r.FormFile("logo")

	if err != nil {
		responseJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	description := r.FormValue("description")
	content := r.FormValue("content")
	isActive := r.FormValue("is_active")
	openedAt := r.FormValue("opened_at")
	closedAt := r.FormValue("closed_at")

	isActiveBool, _ := strconv.ParseBool(isActive)

	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Println(err.Error())
	}

	payloads := struct {
		Title       string `validate:"required"`
		Logo        string `validate:"required"`
		Description string `validate:"required"`
		Content     string `validate:"required"`
		IsActive    bool   `validate:"required"`
		OpenedAt    string `validate:"required"`
		ClosedAt    string `validate:"required"`
	}{
		Title:       title,
		Logo:        header.Filename,
		Description: description,
		Content:     content,
		IsActive:    isActiveBool,
		OpenedAt:    openedAt,
		ClosedAt:    closedAt,
	}

	// validate field
	validate = validator.New()
	err = validate.Struct(&payloads)
	if err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, fmt.Sprintf("Bidang dengan nama %s ini dibutuhkan.", err.Field()))
		}
		responseJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":   "error",
			"messsage": "Harap perbaiki validasi form.",
			"errors":   errorMessages,
		})
		return
	}

	// parsedOpenedAt, err := time.Parse("15:04", openedAt)
	// if err != nil {
	// 	responseJSON(w, http.StatusInternalServerError, map[string]interface{}{
	// 		"status": "error",
	// 		"message": err.Error(),
	// 	})
	// 	return
	// }

	parsedOpenedAtWithLocation, err := time.ParseInLocation("15:04", openedAt, location)
	if err != nil {
		responseJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	parsedClosedAtWithLocation, err := time.ParseInLocation("15:04", closedAt, location)
	if err != nil {
		responseJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// try to save to database model
	course := models.Course{
		Title:       payloads.Title,
		Description: payloads.Description,
		Content:     payloads.Content,
		IsActive:    payloads.IsActive,
		OpenedAt:    parsedOpenedAtWithLocation,
		ClosedAt:    parsedClosedAtWithLocation,
	}

	// check extension
	if !checkExtension(filepath.Ext(header.Filename)) {
		responseJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Harap upload gambar dengan extension JPG|PNG saja.",
		})
		return
	}

	// try to save
	err = models.CreateNewCourse(&course)
	if err != nil {
		responseJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// filename
	filename := fmt.Sprintf("/media/courses/%d/%s", course.ID, header.Filename)

	dir, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create dir
	_ = os.MkdirAll(fmt.Sprintf("./media/courses/%d", course.ID), 0777)

	fileLocation := filepath.Join(dir, filename)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		responseJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	defer targetFile.Close()

	// upload
	if _, err := io.Copy(targetFile, file); err != nil {
		responseJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// update course logo
	course.Logo = filename
	models.DB().Save(&course)

	responseJSON(w, http.StatusCreated, course)
}
