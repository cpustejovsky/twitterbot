package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {

	//avoid using DefaultServeMux to prevent compromised pckgs from exposing malicious handlers to the web
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(home))
	mux.Get("/run-twitter-bot", http.HandlerFunc(app.handleSendEmail))

	return mux
}
