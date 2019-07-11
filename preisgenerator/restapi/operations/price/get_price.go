// Code generated by go-swagger; DO NOT EDIT.

package price

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
	strfmt "github.com/go-openapi/strfmt"
	swag "github.com/go-openapi/swag"
)

// GetPriceHandlerFunc turns a function with the right signature into a get price handler
type GetPriceHandlerFunc func(GetPriceParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetPriceHandlerFunc) Handle(params GetPriceParams) middleware.Responder {
	return fn(params)
}

// GetPriceHandler interface for that can handle valid get price params
type GetPriceHandler interface {
	Handle(GetPriceParams) middleware.Responder
}

// NewGetPrice creates a new http.Handler for the get price operation
func NewGetPrice(ctx *middleware.Context, handler GetPriceHandler) *GetPrice {
	return &GetPrice{Context: ctx, Handler: handler}
}

/*GetPrice swagger:route GET /price price getPrice

Get price

Get a random price

*/
type GetPrice struct {
	Context *middleware.Context
	Handler GetPriceHandler
}

func (o *GetPrice) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetPriceParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// GetPriceOKBody get price o k body
// swagger:model GetPriceOKBody
type GetPriceOKBody struct {

	// price
	Price float32 `json:"price,omitempty"`
}

// Validate validates this get price o k body
func (o *GetPriceOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetPriceOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetPriceOKBody) UnmarshalBinary(b []byte) error {
	var res GetPriceOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}