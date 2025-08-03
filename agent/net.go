package main

import (
	"common/profiles"
	"io"
)

func send(agentId *string, profile *profiles.Profile, data *string) (err error) {
	reqBody, err := profiles.OperateData(&profile.Outcall.Payload.Operations, data)
	if err != nil {
		return err
	}
	processedAgentId, err := profiles.OperateData(&profile.Outcall.ID.Operations, agentId)
	if err != nil {
		return err
	}

	client := profile.GetHttpClient()

	req, err := profile.GetRequestOutcall(&processedAgentId, []byte(reqBody))
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

func receive(agentId *string, profile *profiles.Profile) (instruction string, err error) {
	processedAgentId, err := profiles.OperateData(&profile.Outcall.ID.Operations, agentId)
	if err != nil {
		return "", err
	}

	client := profile.GetHttpClient()

	req, err := profile.GetRequestIncall(&processedAgentId, nil)
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

	bodyStr := string(body)

	return profiles.OperateData(&profile.Incall.Payload.Operations, &bodyStr)
}
