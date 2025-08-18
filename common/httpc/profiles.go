package httpc

import (
	"bytes"
	"encoding/json"
	"net/http"

	profiles "github.com/zarkones/ControlPROFILE"
)

func InsertProfile(name, description string, profile profiles.Profile) (err error) {
	serializedProfile, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	bodyCtx := map[string]any{
		"Name":              name,
		"Description":       description,
		"SerializedProfile": string(serializedProfile),
	}

	serializedBodyCtx, err := json.Marshal(bodyCtx)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, BaseURL+"/v1/profiles", bytes.NewBuffer([]byte(serializedBodyCtx)))
	if err != nil {
		return err
	}

	setAuthHeader(req)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return ErrUnexpectedStatusCode
	}
	return nil
}
