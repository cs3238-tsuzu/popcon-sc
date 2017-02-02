package swagger

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetContestsCidJoinUIDParams creates a new GetContestsCidJoinUIDParams object
// with the default values initialized.
func NewGetContestsCidJoinUIDParams() GetContestsCidJoinUIDParams {
	var ()
	return GetContestsCidJoinUIDParams{}
}

// GetContestsCidJoinUIDParams contains all the bound params for the get contests cid join UID operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetContestsCidJoinUID
type GetContestsCidJoinUIDParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request

	/*Contest ID
	  Required: true
	  In: path
	*/
	Cid int64
	/*User ID
	  Required: true
	  In: path
	*/
	UID int64
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *GetContestsCidJoinUIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	rCid, rhkCid, _ := route.Params.GetOK("cid")
	if err := o.bindCid(rCid, rhkCid, route.Formats); err != nil {
		res = append(res, err)
	}

	rUID, rhkUID, _ := route.Params.GetOK("uid")
	if err := o.bindUID(rUID, rhkUID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetContestsCidJoinUIDParams) bindCid(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("cid", "path", "int64", raw)
	}
	o.Cid = value

	return nil
}

func (o *GetContestsCidJoinUIDParams) bindUID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("uid", "path", "int64", raw)
	}
	o.UID = value

	return nil
}
