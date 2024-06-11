package handlers

import (
	"net/http"
	"shortier/db"

	"github.com/gin-gonic/gin"
)

func (h BaseHandler) InsertLink(c *gin.Context) {
	var link *db.Link
	if err := c.BindJSON(&link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.db.CreateHash(*link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}
