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
	grantTypeParam   = "grant_type"
	scopeParam       = "scope"
	audienceParam    = "audience"
	subjectParam     = "subject"
	tokenFormatParam = "token_format"
	defaultAudience  = "default"
	defaultSubject   = "test"
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
	grant := c.PostForm(grantTypeParam)
	scope := c.PostForm(scopeParam)

	tokenFormat := c.Query(tokenFormatParam)
	if tokenFormat == "" {
		tokenFormat = c.PostForm(tokenFormatParam)
	}

	switch grant {
	default:
		fallthrough
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
			audience, audProvided := c.GetPostForm(audienceParam)
			if !audProvided {
				audience = defaultAudience // for backward compatibility
			}

			subject, subProvided := c.GetPostForm(subjectParam)
			if !subProvided {
				subject = defaultSubject // for backward compatibility
			}

			rsaJwt, err := h.NewRSAJWT(scope, audience, subject)
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

func (h Handler) NewRSAJWT(scope string, audience string, subject string) ([]byte, error) {
	builder := jwt.NewBuilder().
		Issuer(h.IssuerURL).
		IssuedAt(time.Now()).
		NotBefore(time.Now().Add(-1 * time.Hour)).
		Expiration(time.Now().Add(1 * time.Hour))

	if subject != "" {
		builder.Subject(subject)
	}
	if audience != "" {
		builder.Audience(strings.Split(audience, ","))
	}
	if scope != "" {
		builder.Claim("scope", scope)
	}

	t, err := builder.Build()
	if err != nil {
		return nil, err
	}

	return jwt.Sign(t, jwt.WithKey(jwa.RS256, h.Key))
}
