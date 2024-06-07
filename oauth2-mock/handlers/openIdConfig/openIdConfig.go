package openIdConfig

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	Iss string
}

func (h Handler) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"issuer":                 h.Iss,
		"jwks_uri":               fmt.Sprintf("%s/oauth2/certs", h.Iss),
		"token_endpoint":         fmt.Sprintf("%s/oauth2/token", c.Request.Host),
		"introspection_endpoint": fmt.Sprintf("%s/oauth2/introspect", h.Iss),
	})
}
