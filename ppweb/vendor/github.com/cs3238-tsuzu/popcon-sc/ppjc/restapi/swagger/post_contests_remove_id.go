package swagger

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// PostContestsRemoveIDHandlerFunc turns a function with the right signature into a post contests remove ID handler
type PostContestsRemoveIDHandlerFunc func(PostContestsRemoveIDParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PostContestsRemoveIDHandlerFunc) Handle(params PostContestsRemoveIDParams) middleware.Responder {
	return fn(params)
}

// PostContestsRemoveIDHandler interface for that can handle valid post contests remove ID params
type PostContestsRemoveIDHandler interface {
	Handle(PostContestsRemoveIDParams) middleware.Responder
}

// NewPostContestsRemoveID creates a new http.Handler for the post contests remove ID operation
func NewPostContestsRemoveID(ctx *middleware.Context, handler PostContestsRemoveIDHandler) *PostContestsRemoveID {
	return &PostContestsRemoveID{Context: ctx, Handler: handler}
}

/*PostContestsRemoveID swagger:route POST /contests/remove/{id} postContestsRemoveId

Remove the ranking for the contest


*/
type PostContestsRemoveID struct {
	Context *middleware.Context
	Handler PostContestsRemoveIDHandler
}

func (o *PostContestsRemoveID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewPostContestsRemoveIDParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
