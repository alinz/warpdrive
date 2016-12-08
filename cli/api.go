package cli

import (
	"fmt"

	"encoding/json"

	"github.com/pressly/warpdrive/data"
)

type api struct {
	serverAddr string
	session    string
}

func (a *api) validate() error {
	path, err := joinURL(a.serverAddr, "/session/start")
	if err != nil {
		return fmt.Errorf("Server Address '%s' is invalid", a.serverAddr)
	}

	_, err = httpRequest("GET", path, nil, a.session)
	return err
}

func (a *api) login(email, password string) error {
	reqBody := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    email,
		Password: password,
	}

	path, err := joinURL(a.serverAddr, "/session/start")
	if err != nil {
		return fmt.Errorf("Server Address '%s' is invalid", a.serverAddr)
	}

	resp, err := httpRequest("POST", path, reqBody, "")
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("login failed, status code: %d", resp.StatusCode)
	}

	// grab the jwt from resceived cookie from server and assign it to api's session
	cookie := resp.Cookies()[0]
	a.session = cookie.Value

	return nil
}

func (a *api) getCycle(appID, cycleID int64) (*data.Cycle, error) {
	path, err := joinURL(a.serverAddr, fmt.Sprintf("/apps/%d/cycles/%d", appID, cycleID))
	if err != nil {
		return nil, fmt.Errorf("Server Address '%s' is invalid", a.serverAddr)
	}

	resp, err := httpRequest("POST", path, nil, a.session)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("don't access to this app")
	}

	var cycle data.Cycle

	err = json.NewDecoder(resp.Body).Decode(&cycle)
	if err != nil {
		return nil, err
	}

	return &cycle, nil
}

func newAPI(serverAddr string) *api {
	return &api{serverAddr: serverAddr}
}
