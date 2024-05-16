package jwks

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

type Handler struct {
	JWK jwk.Key
}

func (h Handler) Handle(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"keys": []interface{}{
			h.JWK,
		},
	})
}
