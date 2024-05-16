package introspect

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

const (
	active = "active"
	scope  = "scope"
)

type Handler struct {
	ReadToken    uuid.UUID
	NoScopeToken uuid.UUID
}

func NewHandler(readToken uuid.UUID, noScopeToken uuid.UUID) *Handler {
	return &Handler{ReadToken: readToken, NoScopeToken: noScopeToken}
}

func (h Handler) Handle(c *gin.Context) {
	token := c.PostForm("token")

	switch token {
	case h.ReadToken.String():
		c.JSON(http.StatusOK, gin.H{
			active: true,
			scope:  "read",
		})
	case h.NoScopeToken.String():
		c.JSON(http.StatusOK, gin.H{
			active: true,
		})
	default:
		c.Status(http.StatusUnauthorized)
	}
}
