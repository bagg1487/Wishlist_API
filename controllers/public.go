package controllers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"wishlist-api/models"
	"wishlist-api/utils"
	"wishlist-api/database"
)

// GetPublicWishlist godoc
// @Summary Публичный вишлист
// @Tags public
// @Produce json
// @Param token path string true "Public token"
// @Success 200 {object} map[string]interface{}
// @Router /public/{token} [get]
func GetPublicWishlist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	var wishlist models.Wishlist
	db := database.GetDB()
	if err := db.Where("public_token = ?", token).First(&wishlist).Error; err != nil {
		utils.Respond(w, utils.Message(false, "Wishlist not found"))
		return
	}

	var items []models.Item
	db.Where("wishlist_id = ?", wishlist.ID).Find(&items)

	for i := range items {
		var booking models.Booking
		if db.Where("item_id = ?", items[i].ID).First(&booking).Error == nil {
			items[i].IsBooked = true
		}
	}

	resp := utils.Message(true, "success")
	resp["wishlist"] = wishlist
	resp["items"] = items
	utils.Respond(w, resp)
}

// BookItem godoc
// @Summary Забронировать item
// @Tags public
// @Produce json
// @Param token path string true "Public token"
// @Param itemId path int true "Item ID"
// @Param booked_by query string true "Кто бронирует"
// @Success 200 {object} map[string]interface{}
// @Router /public/{token}/book/{itemId} [post]
func BookItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	itemID, _ := strconv.Atoi(vars["itemId"])

	bookedBy := r.URL.Query().Get("booked_by")
	if bookedBy == "" {
		utils.Respond(w, utils.Message(false, "booked_by parameter required"))
		return
	}

	db := database.GetDB()

	var wishlist models.Wishlist
	if err := db.Where("public_token = ?", token).First(&wishlist).Error; err != nil {
		utils.Respond(w, utils.Message(false, "Wishlist not found"))
		return
	}

	var item models.Item
	if err := db.Where("id = ? AND wishlist_id = ?", itemID, wishlist.ID).First(&item).Error; err != nil {
		utils.Respond(w, utils.Message(false, "Item not found"))
		return
	}

	var existingBooking models.Booking
	if err := db.Where("item_id = ?", itemID).First(&existingBooking).Error; err == nil {
		utils.Respond(w, utils.Message(false, "Item already booked"))
		return
	}

	booking := &models.Booking{
		ItemID:   uint(itemID),
		BookedBy: bookedBy,
	}
	db.Create(booking)

	resp := utils.Message(true, "Item booked")
	resp["booking"] = booking
	utils.Respond(w, resp)
}