package profiles

import "encoding/json"

type ReceiveRequest struct {
	SleepMin        int                 `json:"sleepMin"`
	SleepMax        int                 `json:"sleepMax"`
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

type Profile struct {
	// PublicKey PublicKey      `json:"publicKey"`
	Receive ReceiveRequest `json:"receive"`
}

func (p *Profile) Validate() (err error) {
	// TODO
	return nil
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
