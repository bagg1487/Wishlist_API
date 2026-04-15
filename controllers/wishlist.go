package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"wishlist-api/models"
	"wishlist-api/utils"
	"wishlist-api/database"
)

// CreateWishlist godoc
// @Summary Создать вишлист
// @Tags wishlists
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "title, description, event_date"
// @Success 200 {object} map[string]interface{}
// @Router /wishlists [post]

func CreateWishlist(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user").(uint)
	if !ok {
		utils.Respond(w, utils.Message(false, "Unauthorized"))
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		EventDate   string `json:"event_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Respond(w, utils.Message(false, "Invalid request"))
		return
	}

	if req.Title == "" {
		utils.Respond(w, utils.Message(false, "Title required"))
		return
	}

	var eventDate *time.Time
	if req.EventDate != "" {
		t, _ := time.Parse("2006-01-02", req.EventDate)
		eventDate = &t
	}

	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	publicToken := hex.EncodeToString(tokenBytes)

	wishlist := &models.Wishlist{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		EventDate:   eventDate,
		PublicToken: publicToken,
	}

	db := database.GetDB()
	db.Create(wishlist)

	resp := utils.Message(true, "Wishlist created")
	resp["wishlist"] = wishlist
	utils.Respond(w, resp)
}

// GetWishlists godoc
// @Summary Получить все вишлисты
// @Tags wishlists
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /wishlists [get]
func GetWishlists(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user").(uint)
	if !ok {
		utils.Respond(w, utils.Message(false, "Unauthorized"))
		return
	}

	var wishlists []models.Wishlist
	db := database.GetDB()
	db.Where("user_id = ?", userID).Order("created_at DESC").Find(&wishlists)

	resp := utils.Message(true, "success")
	resp["data"] = wishlists
	utils.Respond(w, resp)
}

// GetWishlist godoc
// @Summary Получить один вишлист
// @Tags wishlists
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} map[string]interface{}
// @Router /wishlists/{id} [get]
func GetWishlist(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("user").(uint)
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var wishlist models.Wishlist
	db := database.GetDB()
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&wishlist).Error; err != nil {
		utils.Respond(w, utils.Message(false, "Not found"))
		return
	}

	var items []models.Item
	db.Where("wishlist_id = ?", wishlist.ID).Find(&items)

	resp := utils.Message(true, "success")
	resp["wishlist"] = wishlist
	resp["items"] = items
	utils.Respond(w, resp)
}

// UpdateWishlist godoc
// @Summary Обновить вишлист
// @Tags wishlists
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param request body object true "данные"
// @Success 200 {object} map[string]interface{}
// @Router /wishlists/{id} [put]
func UpdateWishlist(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("user").(uint)
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		EventDate   string `json:"event_date"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	var wishlist models.Wishlist
	db := database.GetDB()
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&wishlist).Error; err != nil {
		utils.Respond(w, utils.Message(false, "Not found"))
		return
	}

	if req.Title != "" {
		wishlist.Title = req.Title
	}
	wishlist.Description = req.Description
	if req.EventDate != "" {
		t, _ := time.Parse("2006-01-02", req.EventDate)
		wishlist.EventDate = &t
	}

	db.Save(&wishlist)

	resp := utils.Message(true, "Updated")
	resp["wishlist"] = wishlist
	utils.Respond(w, resp)
}

// DeleteWishlist godoc
// @Summary Удалить вишлист
// @Tags wishlists
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} map[string]interface{}
// @Router /wishlists/{id} [delete]
func DeleteWishlist(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("user").(uint)
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	db := database.GetDB()
	result := db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Wishlist{})
	if result.RowsAffected == 0 {
		utils.Respond(w, utils.Message(false, "Not found"))
		return
	}

	utils.Respond(w, utils.Message(true, "Deleted"))
}