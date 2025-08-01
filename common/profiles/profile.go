package profiles

import (
	"bytes"
	"common/slices"
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type ProfileRequest struct {
	Timeout         int                 `json:"timeout"`
	Hosts           []string            `json:"hosts"`
	Routes          []string            `json:"routes"`
	Headers         map[string][]string `json:"headers"`
	Methods         []string            `json:"methods"`
	Operations      []string            `json:"operations"`
	PayloadPosition []string            `json:"payloadPosition"`
	UrlParams       map[string][]string `json:"urlParams"`
	Prefix          string              `json:"prefix"`
	Suffix          string              `json:"suffix"`
}

// type PublicKey struct {
// 	Encodings []string `json:"encodings"`
// 	Value     string   `json:"value"`
// }

type Tick struct {
	SleepMin int `json:"sleepMin"`
	SleepMax int `json:"sleepMax"`
}

type Profile struct {
	// PublicKey PublicKey      `json:"publicKey"`
	Incall  ProfileRequest `json:"incall"`
	Outcall ProfileRequest `json:"outcall"`
	Tick    Tick           `json:"tick"`
}

func (p *Profile) Validate() (err error) {
	// TODO
	return nil
}

func (p *Profile) GetTickAmount() (amount time.Duration) {
	if p.Tick.SleepMin == p.Tick.SleepMax {
		return time.Millisecond * time.Duration(p.Tick.SleepMin)
	}
	return time.Millisecond * time.Duration(rand.Intn(p.Tick.SleepMax-p.Tick.SleepMin)+p.Tick.SleepMax)
}

func (p *Profile) GetHttpClient() (client *http.Client) {
	return &http.Client{
		Timeout: time.Millisecond * time.Duration(p.Incall.Timeout),
	}
}

func (p *Profile) GetRequestIncall(body []byte) (req *http.Request, err error) {
	return p.getRequest(&p.Incall, body)
}

func (p *Profile) GetRequestOutcall(body []byte) (req *http.Request, err error) {
	return p.getRequest(&p.Outcall, body)
}

func (p *Profile) getRequest(profileRequest *ProfileRequest, body []byte) (req *http.Request, err error) {
	method := slices.Rand(&profileRequest.Methods)
	host := slices.Rand(&profileRequest.Hosts)
	route := "/" + strings.TrimPrefix(slices.Rand(&profileRequest.Routes), "/")

	if len(profileRequest.UrlParams) != 0 {
		route += "?"
		for name, values := range profileRequest.UrlParams {
			route += name + "=" + slices.Rand(&values) + "&"
		}
		route = strings.TrimSuffix(route, "&")
	}

	var bodyBuffer *bytes.Buffer = bytes.NewBuffer(nil)
	if len(body) != 0 {
		bodyBuffer.Write(body)
	}

	req, err = http.NewRequest(method, host+route, bodyBuffer)
	if err != nil {
		return nil, err
	}

	if len(profileRequest.Headers) != 0 {
		for name, values := range profileRequest.Headers {
			if name == "Host" {
				req.Host = slices.Rand(&values)
			}
			req.Header.Add(name, slices.Rand(&values))
		}
	}

	return req, nil
}

func Load(rawProfile []byte) (profile Profile, err error) {
	if err := json.Unmarshal(rawProfile, &profile); err != nil {
		return profile, nil
	}

	if err := profile.Validate(); err != nil {
		return profile, nil
	}

	return profile, nil
}
