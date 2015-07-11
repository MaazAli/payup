package user

import (
	"net/http"
	"payup/database"

	"github.com/gin-gonic/gin"
)

// Index is used when the user's index is routed to
// this handler will run. Generally, it will
// come with some query parameters like limit and offset
// @returns an array of users
func Index(c *gin.Context) {
	var users []User
	database.DBCon.Limit(c.Param("limit")).Find(&users)

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// Show is used to show one specific user
// @returns a user struct
func Show(c *gin.Context) {
	var user User
	database.DBCon.First(&user, c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Create is used to create one specific user, it'll come with some form data
// @returns the newly created user struct
func Create(c *gin.Context) {

	user := User{
		Name:     c.PostForm("name"),
		Username: c.PostForm("username"),
		Email:    c.PostForm("email"),
	}

	database.DBCon.Create(&user)

	c.JSON(http.StatusOK, gin.H{"user": user})

}

// Update is used to update a specific user, it'll also come with some form data
// @returns a user struct
func Update(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"userUpdate": "someContent"})
}

// Delete is used to delete one specific user with a `id`
func Delete(c *gin.Context) {
	var user User
	database.DBCon.First(&user, c.Param("id"))
	database.DBCon.Delete(&user)

	c.JSON(http.StatusOK, gin.H{"success": true})
}
