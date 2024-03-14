package controllers

import (
	"Aura-Server/initializers"
	"Aura-Server/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRoom(c *gin.Context) {
	// Get the credentials off req body
	var body struct {
		Name string
	}

	// Get user from middleware
	user, _ := c.Get("user")

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	// Craete new room
	room := models.Room{
		Name:      body.Name,
		CreatedBy: user.(models.User).ID,
	}
	result := initializers.DB.Create(&room)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create room",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func UpdateRoom(c *gin.Context) {
	// Get the vars off req body
	var body struct {
		ID   string
		Name string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	result := initializers.DB.Model(&models.Room{}).Where("id = ?", body.ID).
		Updates(models.Room{Name: body.Name})

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update room",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully updated the room"})
}

func DeleteRoom(c *gin.Context) {
	// Get the id off req body
	var body struct {
		ID string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	editResult := initializers.DB.Model(&models.Device{}).Where("room_id = ?", body.ID).
		Updates(map[string]interface{}{"Configured": false, "RoomID": nil})

	deleteResult := initializers.DB.Where("id = ?", body.ID).Delete(&models.Room{})

	if editResult.Error != nil || deleteResult.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete room",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully deleted room"})
}
