package api

import (
	"context"
	"crafting/config"
	"net/http"

	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/http/middleware"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

func contentTypeJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, req)
	})
}

func authentication() func(http.Handler) http.Handler {
	ctx := context.Background()

	introspectionAuth := oauth.ClientIDSecretIntrospectionAuthentication(config.LoadedConfiguration.OidcServerClientId, config.LoadedConfiguration.OidcServerClientSecret)
	zitadelConfig := oauth.WithIntrospection[*oauth.IntrospectionContext](introspectionAuth)
	authZ, err := authorization.New(ctx, zitadel.New(config.LoadedConfiguration.OidcDomain), zitadelConfig)

	if err != nil {
		panic(err)
	}

	mw := middleware.New(authZ)

	return mw.RequireAuthorization()
}
