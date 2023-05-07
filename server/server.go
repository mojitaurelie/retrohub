package server

import (
	"embed"
	_ "embed"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mileusna/useragent"
	"html/template"
	"log"
	"net/http"
	"os"
	"retroHub/data"
	"strings"
)

type injector struct {
	Version string
}

type userAgentInjector struct {
	injector
	UserAgent useragent.UserAgent
}

type providerInjector struct {
	injector
	Provider data.Provider
}

type errorInjector struct {
	injector
	StatusCode int
	StatusText string
}

const version string = "0.1"

//go:embed templates
var templates embed.FS

var (
	provider      data.Provider
	indexTemplate *template.Template
	uaTemplate    *template.Template
	errTemplate   *template.Template

	errLog *log.Logger
)

func init() {
	indexTemplate = template.Must(template.ParseFS(templates, "templates/base.html", "templates/index.html"))
	uaTemplate = template.Must(template.ParseFS(templates, "templates/base.html", "templates/ua.html"))
	errTemplate = template.Must(template.ParseFS(templates, "templates/base.html", "templates/error.html"))

	errLog = log.New(os.Stderr, "", log.LstdFlags)
}

func Serve(contentProvider data.Provider, port uint) error {
	if contentProvider == nil {
		return errors.New("content provider cannot be nil")
	}
	provider = contentProvider

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(panicHandler)

	r.NotFound(notFoundHandler)
	r.Get("/", indexHandler)
	r.Get("/ua", userAgentHandler)

	log.Printf("the server is up and running on %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		return err
	}
	return nil
}

func panicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				errLog.Println("ERROR", r)
				injectedData := errorInjector{
					injector: injector{
						Version: version,
					},
					StatusCode: http.StatusInternalServerError,
					StatusText: "Internal Server Error",
				}
				w.WriteHeader(http.StatusInternalServerError)
				err := errTemplate.ExecuteTemplate(w, "base", injectedData)
				if err != nil {
					errLog.Println("ERROR", err)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	injectedData := providerInjector{
		injector: injector{
			Version: version,
		},
		Provider: provider,
	}
	err := indexTemplate.ExecuteTemplate(w, "base", injectedData)
	if err != nil {
		errLog.Println("ERROR", err)
	}
}

func userAgentHandler(w http.ResponseWriter, r *http.Request) {
	injectedData := userAgentInjector{
		injector: injector{
			Version: version,
		},
	}
	userAgent := strings.TrimSpace(r.UserAgent())
	if len(userAgent) > 0 {
		injectedData.UserAgent = useragent.Parse(userAgent)
	}
	err := uaTemplate.ExecuteTemplate(w, "base", injectedData)
	if err != nil {
		errLog.Println("ERROR", err)
	}
}

func notFoundHandler(w http.ResponseWriter, _ *http.Request) {
	injectedData := errorInjector{
		injector: injector{
			Version: version,
		},
		StatusCode: http.StatusNotFound,
		StatusText: "Not Found",
	}
	w.WriteHeader(http.StatusNotFound)
	err := errTemplate.ExecuteTemplate(w, "base", injectedData)
	if err != nil {
		errLog.Println("ERROR", err)
	}
}
