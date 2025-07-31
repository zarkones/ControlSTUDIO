package main

import (
	"common/profiles"
	"common/slices"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/zarkones/netescape"
)

var (
	ErrUnknownOp = errors.New("unknown operation")
)

func receive(profile *profiles.Profile) (instruction string, err error) {
	client := &http.Client{
		Timeout: time.Millisecond * time.Duration(profile.Receive.Timeout),
	}

	req, err := reqFromProfile(profile)
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

	return processRespBody(profile, &body)
}

func processRespBody(profile *profiles.Profile, body *[]byte) (data string, err error) {
	data = string(*body)

	for _, operation := range profile.Receive.Operations {
		switch operation {
		default:
			return data, ErrUnknownOp
		case "hex":
			decoded, err := hex.DecodeString(data)
			if err != nil {
				return "", err
			}
			data = string(decoded)
		case "base64_std":
			decoded, err := base64.StdEncoding.DecodeString(data)
			if err != nil {
				return "", err
			}
			data = string(decoded)
		case "base64_url":
			decoded, err := base64.URLEncoding.DecodeString(data)
			if err != nil {
				return "", err
			}
			data = string(decoded)
		case "csv":
			decoded, err := netescape.FromCSV(&data)
			if err != nil {
				return "", err
			}
			data = decoded
		}
	}

	return data, nil
}

func reqFromProfile(profile *profiles.Profile) (req *http.Request, err error) {
	method := slices.Rand(&profile.Receive.Methods)
	host := slices.Rand(&profile.Receive.Hosts)
	route := "/" + strings.TrimPrefix(slices.Rand(&profile.Receive.Routes), "/")

	if len(profile.Receive.UrlParams) != 0 {
		route += "?"
		for name, values := range profile.Receive.UrlParams {
			route += name + "=" + slices.Rand(&values) + "&"
		}
		route = strings.TrimSuffix(route, "&")
	}

	req, err = http.NewRequest(method, host+route, nil)
	if err != nil {
		return nil, err
	}

	if len(profile.Receive.Headers) != 0 {
		for name, values := range profile.Receive.Headers {
			if name == "Host" {
				req.Host = slices.Rand(&values)
			}
			req.Header.Add(name, slices.Rand(&values))
		}
	}

	return req, nil
}
