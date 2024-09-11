package controllers

//AdminController

import (
	"backendGoAuth/internal/entities"
	"backendGoAuth/internal/services"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type AdminController struct {
	service services.AdminService
}

func NewAdminController(service services.AdminService) *AdminController {
	return &AdminController{service}
}

//func (c *AdminController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
//	users, err := c.service.GetAllUsers()
//	if err != nil {
//		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	err = json.NewEncoder(w).Encode(users)
//	if err != nil {
//		return
//	}
//}

func (c *AdminController) GetAllUsers(ctx *gin.Context) {
	users, err := c.service.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (c *AdminController) EditUser(w http.ResponseWriter, r *http.Request) {
	var user entities.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := c.service.EditUser(user); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("User updated successfully"))
	if err != nil {
		return
	}
}

func (c *AdminController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := c.service.DeleteUser(userID); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("User deleted successfully"))
	if err != nil {
		return
	}
}
