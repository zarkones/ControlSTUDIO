package main

import (
	"common/profiles"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"

	"github.com/zarkones/netescape"
)

var (
	ErrUnknownOp = errors.New("unknown operation")
)

func send(profile *profiles.Profile, data *string) (err error) {
	body := []byte(*data)
	reqBody, err := processBody(&profile.Outcall, &body)
	if err != nil {
		return err
	}

	client := profile.GetHttpClient()

	req, err := profile.GetRequestOutcall([]byte(reqBody))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	return nil
}

func receive(profile *profiles.Profile) (instruction string, err error) {
	client := profile.GetHttpClient()

	req, err := profile.GetRequestIncall(nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if len(body) == 0 {
		return "", nil
	}

	return processBody(&profile.Incall, &body)
}

func processBody(req *profiles.ProfileRequest, body *[]byte) (data string, err error) {
	data = string(*body)

	for _, operation := range req.Operations {
		switch operation {

		default:
			return data, ErrUnknownOp

		// HEX
		case "hex_d":
			decoded, err := hex.DecodeString(data)
			if err != nil {
				return "", err
			}
			data = string(decoded)
		case "hex":
			data = hex.EncodeToString([]byte(data))

		// BASE64 STD
		case "base64_std_d":
			decoded, err := base64.StdEncoding.DecodeString(data)
			if err != nil {
				return "", err
			}
			data = string(decoded)
		case "base64_std":
			data = base64.StdEncoding.EncodeToString([]byte(data))

		// BASE64 URL
		case "base64_url_d":
			decoded, err := base64.URLEncoding.DecodeString(data)
			if err != nil {
				return "", err
			}
			data = string(decoded)
		case "base64_url":
			data = base64.URLEncoding.EncodeToString([]byte(data))

		// CSV
		case "csv_d":
			decoded, err := netescape.FromCSV(&data)
			if err != nil {
				return "", err
			}
			data = decoded
		case "csv":
			decoded, err := netescape.ToCsv(&data)
			if err != nil {
				return "", err
			}
			data = decoded

		}
	}

	return data, nil
}
