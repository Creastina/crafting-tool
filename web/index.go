package web

import (
	"crafting/config"
	"embed"
	"html/template"
	fs2 "io/fs"
	"net/http"
	"os"
)

//go:embed tmpl
var tmplFs embed.FS

func indexPage(w http.ResponseWriter, r *http.Request) {
	var fs fs2.FS
	fs = tmplFs
	if config.LoadedConfiguration.Env == "dev" {
		fs = os.DirFS("web")
	}
	t, err := template.New("content").ParseFS(fs, "tmpl/index.gohtml")
	if err == nil {
		t.Execute(w, struct {
			OidcFrontendClientId string
			OidcDomain           string
			ServerUrl            string
		}{
			OidcFrontendClientId: os.Getenv("OIDC_FRONTEND_CLIENT_ID"),
			OidcDomain:           os.Getenv("OIDC_DOMAIN"),
			ServerUrl:            os.Getenv("SERVER_URL"),
		})
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
