package main

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/google/uuid"
	"github.com/kyma-project/networking-dev-tools/oauth2-mock/handlers/introspect"
	"github.com/kyma-project/networking-dev-tools/oauth2-mock/handlers/jwks"
	"github.com/kyma-project/networking-dev-tools/oauth2-mock/handlers/openIdConfig"
	"github.com/kyma-project/networking-dev-tools/oauth2-mock/handlers/token"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	readToken := uuid.New()
	noScopeToken := uuid.New()
	readWriteToken := uuid.New()

	issuerUrl, ok := os.LookupEnv("iss")
	if !ok {
		issuerUrl = "http://mock-oauth2.kyma-system.svc.cluster.local"
	}

	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	jwkFromKey, err := jwk.FromRaw(rsaKey)
	if err != nil {
		panic(err)
	}

	err = jwkFromKey.Set("kid", "1")
	if err != nil {
		panic(err)
	}

	err = jwkFromKey.Set("use", "sig")
	if err != nil {
		panic(err)
	}

	err = jwkFromKey.Set("alg", "RS256")
	if err != nil {
		panic(err)
	}

	publicKey, err := jwkFromKey.PublicKey()
	if err != nil {
		panic(err)
	}

	tokenHandler := token.NewHandler(readToken, readWriteToken, noScopeToken, issuerUrl, jwkFromKey)
	introspectHandler := introspect.NewHandler(readToken, readWriteToken, noScopeToken)
	jwksHandler := jwks.Handler{JWK: publicKey}
	openIdHandler := openIdConfig.Handler{Iss: issuerUrl}

	r.POST("/oauth2/token", tokenHandler.Handle)
	r.POST("/oauth2/introspect", introspectHandler.Handle)
	r.GET("/oauth2/certs", jwksHandler.Handle)
	r.GET("/.well-known/openid-configuration", openIdHandler.Handle)

	err = r.Run()
	if err != nil {
		return
	} // listen and serve on :8080
}
