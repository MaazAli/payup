package controllers

import (
	"net/http"
	"payup/database"
	"payup/models"

	"github.com/gin-gonic/gin"
)

// GroupIndex When the group's index is routed to
// this handler will run. Generally, it will
// come with some query parameters like limit and offset
// @returns an array of group structs
func GroupIndex(c *gin.Context) {
	var groups []models.Group
	database.DBCon.Limit(c.Param("limit")).Find(&groups)

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

// GroupShow is used to show one specific group, returns a group struct
// @returns a group struct
func GroupShow(c *gin.Context) {
	var group models.Group
	var users []models.User

	database.DBCon.First(&group, c.Param("id"))
	database.DBCon.Model(&group).Related(&users, "Users")
	group.Users = users
	// database.DBCon.Find(&users)

	// database.DBCon.Model(&group).Association("Users").Append(users)

	// hd := hashids.NewData()
	// hd.Salt = "9398dfajsie288sawiehg"
	// hd.MinLength = 6
	// h := hashids.NewWithData(hd)
	//
	// // Decode
	// e := h.Decode(group.HashID)
	// fmt.Println(e)

	c.JSON(http.StatusOK, gin.H{"group": group})
}

// GroupCreate is used to create one specific group, it'll come with some form data
// @returns a group struct
func GroupCreate(c *gin.Context) {

	var group models.Group
	c.Bind(&group)
	c.JSON(200, group)
	// database.DBCon.Create(&group)

	// c.JSON(http.StatusOK, gin.H{"group": c.PostForm("group")})
}

// GroupUpdate is used to update a specific group, it'll also come with some form data'
// @returns a group struct
func GroupUpdate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"groupUpdate": "someContent"})
}

// GroupDelete is used to delete one specific group with a `id`
func GroupDelete(c *gin.Context) {
	var group models.Group
	database.DBCon.First(&group, c.Param("id"))
	database.DBCon.Delete(&group)

	c.JSON(http.StatusOK, gin.H{"success": true})
}
