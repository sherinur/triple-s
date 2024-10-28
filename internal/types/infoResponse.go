package types

import "encoding/xml"

type InfoResponse struct {
	XMLName xml.Name `xml:"InfoResponse"`
	Info    string   `xml:"Info"`
}

func NewInfoResponse(infoMsg string) *InfoResponse {
	return &InfoResponse{
		Info: infoMsg,
	}
}
