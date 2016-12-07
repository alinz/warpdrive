package cli

import (
	"fmt"
	"path"

	"encoding/json"

	"github.com/pressly/warpdrive/data"
)

type api struct {
	serverAddr string
	session    string
}

func (a *api) login(email, password string) error {
	reqBody := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    email,
		Password: password,
	}

	resp, err := httpRequest("POST", path.Join(a.serverAddr, "/session/start"), reqBody, "")
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
	apiPath := fmt.Sprintf("/apps/%d/cycles/%d", appID, cycleID)

	resp, err := httpRequest("POST", path.Join(a.serverAddr, apiPath), nil, a.session)
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
