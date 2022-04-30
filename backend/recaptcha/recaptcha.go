package recaptcha

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/eyalch/kurzr/backend/core"
)

const verifyURL = "https://www.google.com/recaptcha/api/siteverify"

type reCAPTCHAVerifier struct {
	secret         string
	scoreThreshold float32
}

func NewReCAPTCHAVerifier(secret string, scoreThreshold float32) core.ReCAPTCHAVerifier {
	return &reCAPTCHAVerifier{secret, scoreThreshold}
}

type verifyResponse struct {
	Success            bool     `json:"success"`
	Score              float32  `json:"score"`
	Action             string   `json:"action"`
	ChallengeTimestamp string   `json:"challenge_ts"`
	Hostname           string   `json:"hostname"`
	ErrorCodes         []string `json:"error-codes"`
}

func (rv *reCAPTCHAVerifier) Verify(response string, action string) (bool, error) {
	req, _ := http.NewRequest(http.MethodPost, verifyURL, nil)

	q := req.URL.Query()
	q.Add("secret", rv.secret)
	q.Add("response", response)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, errors.Wrap(err, "could not send request")
	}
	defer resp.Body.Close()

	var r verifyResponse
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return false, errors.Wrap(err, "could not unmarshal response body")
	}

	return r.Success && r.Action == action && r.Score >= rv.scoreThreshold, nil
}
