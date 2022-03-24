package auth

import (
	"encoding/json"
	"io"
)

type AuthenticateRequest struct {
	Username string
	Password string
}

func (r *AuthenticateRequest) Parse(body []byte) error {
	return json.Unmarshal(body, &r)
}

type AuthenticateResponse struct {
	Token string
}

func (r *AuthenticateResponse) Render(writer io.Writer) error {
	response, err := json.Marshal(r)
	if err != nil {
		return err
	}

	_, err = writer.Write(response)
	return err
}
