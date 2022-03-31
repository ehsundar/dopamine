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

type GrantPermissionRequest struct {
	Username   string
	Permission string
}

func (r *GrantPermissionRequest) Parse(reader io.Reader) error {
	d := json.NewDecoder(reader)
	return d.Decode(&r)
}

type DropPermissionRequest struct {
	Username   string
	Permission string
}

func (r *DropPermissionRequest) Parse(reader io.Reader) error {
	d := json.NewDecoder(reader)
	return d.Decode(&r)
}
