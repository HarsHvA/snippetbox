package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *applcation) routes() http.Handler {
	standardMiddleWare := alice.New(app.recoverPanic, app.logRequests, secureHeaders)

	dynamicMiddleWare := alice.New(app.session.Enable, noSurf, app.authenticate)

	mux := pat.New()
	mux.Get("/", dynamicMiddleWare.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleWare.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleWare.Append(app.requireAuthentication).ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleWare.ThenFunc(app.showSnippet))

	mux.Get("/user/signup", dynamicMiddleWare.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleWare.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleWare.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleWare.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleWare.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleWare.Then(mux)
}
