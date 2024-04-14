package web

import (
	"embed"
	"log"
	"net/http"
	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/vaporvee/acecore/custom"
)

// Embed the HTML file into the binary
//
//go:embed html/privacy.html
var privacyHTML embed.FS

//go:embed html/tos.html
var tosHTML embed.FS

func handleHTML(w http.ResponseWriter, embed embed.FS, path string) {
	tmpl, err := template.ParseFS(embed, path)
	if err != nil {
		logrus.Error(err)
		return
	}
	tmpl.Execute(w, nil)
}

func HostRoutes(botID string) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered from panic: %v", r)
			go HostRoutes(botID)
		}
	}()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, custom.Gh_url, http.StatusMovedPermanently)
	})
	http.HandleFunc("/invite", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://discord.com/oauth2/authorize?client_id="+botID, http.StatusMovedPermanently)
	})
	http.HandleFunc("/privacy", func(w http.ResponseWriter, r *http.Request) {
		handleHTML(w, privacyHTML, "html/privacy.html")
	})
	http.HandleFunc("/tos", func(w http.ResponseWriter, r *http.Request) {
		handleHTML(w, tosHTML, "html/tos.html")
	})

	server := &http.Server{
		Addr:     ":443",
		Handler:  nil,
		ErrorLog: log.New(nil, "", 0),
	}
	logrus.Info("Starting server for html routes on :443...")
	if err := server.ListenAndServeTLS("./web/cert.pem", "./web/key.pem"); err != nil {
		logrus.Warnf("Couldn't start server for html routes: %v\n", err)
	}
}
