// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"

	"github.com/ildomm/eskalationszeit/preisgenerator/restapi/operations"
	"github.com/ildomm/eskalationszeit/preisgenerator/restapi/operations/price"
)

//go:generate swagger generate server --target ..\..\preisgenerator --name Preis --spec ..\spec\spec.json

func configureFlags(api *operations.PreisAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.PreisAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()
	
	api.PriceGetPriceHandler = price.GetPriceHandlerFunc(func(params price.GetPriceParams) middleware.Responder {
		priceSuggestion := 10 + rand.Float32()*(99-10)
		body := price.GetPriceOKBody{Price:priceSuggestion}
		return price.NewGetPriceOK().WithPayload(&body)
	})

	api.OptionsAllowHandler = operations.OptionsAllowHandlerFunc(func(params operations.OptionsAllowParams) middleware.Responder {
		return operations.NewOptionsAllowOK()
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func xsetupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}


func setupGlobalMiddleware(handler http.Handler) http.Handler {
	log.Println("Setup GlobalMiddleware")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Token")

		dumpRequest, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("RemoteAddr: "+r.RemoteAddr+", Request: %q", dumpRequest)
		}

		handler.ServeHTTP(w, r)
	})
}
