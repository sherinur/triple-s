package types

import "encoding/xml"

type ErrorResponse struct {
	XMLName xml.Name `xml:"ErrorResponse"`
	Error   string   `xml:"Error"`
	Message string   `xml:"Message"`
}

func NewErrorResponse(errorMsg, message string) *ErrorResponse {
	return &ErrorResponse{
		Error:   errorMsg,
		Message: message,
	}
}
