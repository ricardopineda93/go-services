package account

import (
	"context"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

/*
This can be thought of as the Transport Layer of our service. It deals with the HTTP
specific functionality of directing HTTP endpoint traffic to our Endpoints (which then
safely expose the Service methods according to the Transport type, in this case HTTP.)

In other words, this maps the HTTP-world stuff to a corresponding Endpoint, where
the endpoint can then take the next steps to do some safe "middleman" work to safely "transform"
the Transport request into a Service-friendly format to then safely invoke the Service's methods
*/

// Factory function for creating an HTTP server that will specifically handle the HTTP
// traffic for our Account service. We map HTTP methods and routes to the respective
// Endpoint, where the Endpoint will then broker that HTTP request to the Service
func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	// Init the router
	router := mux.NewRouter()
	// Have the router use the middleware we defined, in this case it simply
	// adds the content-type:application/json header to each of our responses.
	router.Use(commonMiddleware)

	// Instead of passing in the Endpoint directly, we instead
	// do a bit of functional programming-esque stuff where we pass in another function to handler
	// that itself consumes ANOTHER function (the actual endpoint which is itself a function that returns a function)
	router.Methods("POST").Path("/login/org/{org_id}").Handler(
		httptransport.NewServer(
			endpoints.LoginUser,
			DecodeLoginReq,
			EncodeResponse,
		))

	// TODO: Subrouting

	router.Methods("POST").Path("/users/org/{org_id}").Handler(
		httptransport.NewServer( // This "brokers" the Transport and Endpoint layers, allowing us a way to transform the request from on layer to the next
			endpoints.CreateUser, // The endpoint itself
			DecodeCreateUserReq,  // How we want to "decode" the request, i.e. take the HTTP request and cast it into something the Endpoint/service can use
			EncodeResponse,       // How we want to "encode" the resulting response our Endpoint returns (this case as JSON)
		))

	router.Methods("GET").Path("/users/{id}").Handler(httptransport.NewServer(endpoints.GetUser,
		DecodeGetUserReq,
		EncodeResponse,
	))

	router.Methods("PUT").Path("/users/{id}/profile").Handler(httptransport.NewServer(endpoints.UpdateUserProfile,
		DecodeUpdateUserProfileReq,
		EncodeResponse,
	))

	router.Methods("POST").Path("/orgs").Handler(
		httptransport.NewServer(
			endpoints.CreateOrg,
			DecodeCreateOrgReq,
			EncodeResponse,
		))

	return router
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, req *http.Request) {
		respWriter.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(respWriter, req)
	})
}
