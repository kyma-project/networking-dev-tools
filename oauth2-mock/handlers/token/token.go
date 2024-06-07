package token

import (
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"strings"

	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/google/uuid"
	"net/http"
	"time"
)

const (
	grantType = "grant_type"
	scope     = "scope"
	audience  = "audience"
)

type Handler struct {
	ReadToken      uuid.UUID
	NoScopeToken   uuid.UUID
	ReadWriteToken uuid.UUID
	IssuerURL      string
	Key            jwk.Key
}

func NewHandler(readToken uuid.UUID, readWriteToken, noScopeToken uuid.UUID, issuer string, key jwk.Key) *Handler {
	return &Handler{ReadToken: readToken, ReadWriteToken: readWriteToken, NoScopeToken: noScopeToken, IssuerURL: issuer, Key: key}
}

func (h Handler) Handle(c *gin.Context) {
	grant := c.PostForm(grantType)
	scope := c.PostForm(scope)

	tokenFormat := c.Query("token_format")
	if tokenFormat == "" {
		tokenFormat = c.PostForm("token_format")
	}

	switch grant {
	case "client_credentials":
		switch tokenFormat {
		case "opaque":
			switch scope {
			default:
				c.JSON(http.StatusOK, gin.H{
					"access_token": h.ReadWriteToken,
					"token_type":   "Bearer",
					"expires_in":   3600,
				})
			case "read":
				c.JSON(http.StatusOK, gin.H{
					"access_token": h.ReadToken,
					"token_type":   "Bearer",
					"expires_in":   3600,
				})
			case "":
				c.JSON(http.StatusOK, gin.H{
					"access_token": h.NoScopeToken,
					"token_type":   "Bearer",
					"expires_in":   3600,
				})
			}
		default:
			fallthrough
		case "jwt":
			aud := c.PostForm(audience)
			rsaJwt, err := h.NewRSAJWT(scope, aud)
			if err != nil {
				_ = c.Error(err)
				c.Status(http.StatusInternalServerError)
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"access_token": string(rsaJwt),
				"token_type":   "Bearer",
				"expires_in":   3600,
			})
		}
	}
}

func (h Handler) NewRSAJWT(scp string, aud string) ([]byte, error) {
	builder := jwt.NewBuilder().Issuer(h.IssuerURL).NotBefore(time.Now()).IssuedAt(time.Now()).Expiration(time.Now().Add(1 * time.Hour))

	if aud == "" {
		aud = "default"
	}

	builder.Subject("test")
	builder.Audience(strings.Split(aud, ","))

	if scp != "" {
		builder.Claim("scope", scp)
	}

	t, err := builder.Build()
	if err != nil {
		return nil, err
	}

	return jwt.Sign(t, jwt.WithKey(jwa.RS256, h.Key))
}
