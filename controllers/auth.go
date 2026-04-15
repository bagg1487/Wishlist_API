package controllers

import (
	"encoding/json"
	"net/http"

	"wishlist-api/models"
	"wishlist-api/utils"
	"wishlist-api/database"
)

// Register godoc
// @Summary Регистрация пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object true "email и password"
// @Success 200 {object} map[string]interface{}
// @Router /auth/register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Respond(w, utils.Message(false, "Invalid request"))
		return
	}

	if req.Email == "" || req.Password == "" {
		utils.Respond(w, utils.Message(false, "Email and password required"))
		return
	}

	hashed, _ := utils.HashPassword(req.Password)
	user := &models.User{Email: req.Email, PasswordHash: hashed}

	db := database.GetDB()
	if err := db.Create(user).Error; err != nil {
		utils.Respond(w, utils.Message(false, "Email already exists"))
		return
	}

	token := utils.GenerateJWT(user.ID)
	resp := utils.Message(true, "Registered")
	resp["token"] = token
	resp["user"] = user
	utils.Respond(w, resp)
}

// Login godoc
// @Summary Логин пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object true "email и password"
// @Success 200 {object} map[string]interface{}
// @Router /auth/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Respond(w, utils.Message(false, "Invalid request"))
		return
	}

	var user models.User
	db := database.GetDB()
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.Respond(w, utils.Message(false, "Invalid credentials"))
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		utils.Respond(w, utils.Message(false, "Invalid credentials"))
		return
	}

	token := utils.GenerateJWT(user.ID)
	resp := utils.Message(true, "Logged in")
	resp["token"] = token
	resp["user"] = user
	utils.Respond(w, resp)
}