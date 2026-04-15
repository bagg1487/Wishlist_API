package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"wishlist-api/models"
	"wishlist-api/utils"
	"wishlist-api/database"
)

// CreateItem godoc
// @Summary Создать item
// @Tags items
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param wishlistId path int true "Wishlist ID"
// @Param request body object true "данные"
// @Success 200 {object} map[string]interface{}
// @Router /wishlists/{wishlistId}/items [post]
func CreateItem(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("user").(uint)
	vars := mux.Vars(r)
	wishlistID, _ := strconv.Atoi(vars["wishlistId"])

	var wishlist models.Wishlist
	db := database.GetDB()
	if err := db.Where("id = ? AND user_id = ?", wishlistID, userID).First(&wishlist).Error; err != nil {
		utils.Respond(w, utils.Message(false, "Wishlist not found"))
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		URL         string `json:"url"`
		Priority    int    `json:"priority"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if req.Name == "" {
		utils.Respond(w, utils.Message(false, "Name required"))
		return
	}
	if req.Priority < 1 || req.Priority > 5 {
		req.Priority = 1
	}

	item := &models.Item{
		WishlistID:  uint(wishlistID),
		Name:        req.Name,
		Description: req.Description,
		URL:         req.URL,
		Priority:    req.Priority,
	}

	db.Create(item)

	resp := utils.Message(true, "Item created")
	resp["item"] = item
	utils.Respond(w, resp)
}

// GetItems godoc
// @Summary Получить items
// @Tags items
// @Security BearerAuth
// @Produce json
// @Param wishlistId path int true "Wishlist ID"
// @Success 200 {object} map[string]interface{}
// @Router /wishlists/{wishlistId}/items [get]
func GetItems(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("user").(uint)
	vars := mux.Vars(r)
	wishlistID, _ := strconv.Atoi(vars["wishlistId"])

	var wishlist models.Wishlist
	db := database.GetDB()
	if err := db.Where("id = ? AND user_id = ?", wishlistID, userID).First(&wishlist).Error; err != nil {
		utils.Respond(w, utils.Message(false, "Wishlist not found"))
		return
	}

	var items []models.Item
	db.Where("wishlist_id = ?", wishlistID).Order("priority ASC").Find(&items)

	resp := utils.Message(true, "success")
	resp["data"] = items
	utils.Respond(w, resp)
}

// UpdateItem godoc
// @Summary Обновить item
// @Tags items
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Param request body object true "данные"
// @Success 200 {object} map[string]interface{}
// @Router /items/{id} [put]
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("user").(uint)
	vars := mux.Vars(r)
	itemID, _ := strconv.Atoi(vars["id"])

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		URL         string `json:"url"`
		Priority    int    `json:"priority"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	var item models.Item
	db := database.GetDB()
	if err := db.Where("id = ?", itemID).First(&item).Error; err != nil {
		utils.Respond(w, utils.Message(false, "Item not found"))
		return
	}

	var wishlist models.Wishlist
	if err := db.Where("id = ? AND user_id = ?", item.WishlistID, userID).First(&wishlist).Error; err != nil {
		utils.Respond(w, utils.Message(false, "Unauthorized"))
		return
	}

	if req.Name != "" {
		item.Name = req.Name
	}
	item.Description = req.Description
	item.URL = req.URL
	if req.Priority >= 1 && req.Priority <= 5 {
		item.Priority = req.Priority
	}

	db.Save(&item)

	resp := utils.Message(true, "Item updated")
	resp["item"] = item
	utils.Respond(w, resp)
}

// DeleteItem godoc
// @Summary Удалить item
// @Tags items
// @Security BearerAuth
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} map[string]interface{}
// @Router /items/{id} [delete]
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("user").(uint)
	vars := mux.Vars(r)
	itemID, _ := strconv.Atoi(vars["id"])

	var item models.Item
	db := database.GetDB()
	if err := db.Where("id = ?", itemID).First(&item).Error; err != nil {
		utils.Respond(w, utils.Message(false, "Item not found"))
		return
	}

	var wishlist models.Wishlist
	if err := db.Where("id = ? AND user_id = ?", item.WishlistID, userID).First(&wishlist).Error; err != nil {
		utils.Respond(w, utils.Message(false, "Unauthorized"))
		return
	}

	db.Delete(&item)

	utils.Respond(w, utils.Message(true, "Item deleted"))
}