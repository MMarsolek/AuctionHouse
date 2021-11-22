package controller

//go:generate swagger generate spec -m -w ../.. -o ../../swagger.yaml

// ----- Start Documentation Generation Types --------------

// Response contains no body.
// swagger:response noBody
type noBodyDoc struct{}

// Response when there is an error with the request.
// swagger:response errorMessage
type errorMessageDoc struct {
	// In: body
	Body struct {
		// The human readable message to the client.
		Message string `json:"message"`
	}
}

// ----- End Documentation Generation Types --------------
