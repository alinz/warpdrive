package cli

import (
	"fmt"
	"io"

	"encoding/json"

	"github.com/pressly/warpdrive/data"
)

type api struct {
	serverAddr string
	session    string
}

func (a *api) makePath(path string, args ...interface{}) (string, error) {
	path, err := joinURL(a.serverAddr, fmt.Sprintf(path, args...))
	if err != nil {
		return "", fmt.Errorf("Server Address '%s' is invalid", a.serverAddr)
	}
	return path, nil
}

func (a *api) validate() error {
	path, err := a.makePath("/session")
	if err != nil {
		return err
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

	path, err := a.makePath("/session/start")
	if err != nil {
		return err
	}

	resp, err := httpRequest("POST", path, reqBody, "")
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return parseErrorMessage(resp.Body)
	}

	// grab the jwt from resceived cookie from server and assign it to api's session
	cookie := resp.Cookies()[0]
	a.session = cookie.Value

	return nil
}

func (a *api) getCycle(appID, cycleID int64) (*data.Cycle, error) {
	path, err := a.makePath("/apps/%d/cycles/%d", appID, cycleID)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest("GET", path, nil, a.session)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, parseErrorMessage(resp.Body)
	}

	var cycle data.Cycle

	err = json.NewDecoder(resp.Body).Decode(&cycle)
	if err != nil {
		return nil, err
	}

	return &cycle, nil
}

func (a *api) getCycleByName(appID int64, cycleName string) (*data.Cycle, error) {
	path, err := a.makePath("/apps/%d/cycles?name=%s", appID, cycleName)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest("GET", path, nil, a.session)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, parseErrorMessage(resp.Body)
	}

	var cycles []*data.Cycle

	err = json.NewDecoder(resp.Body).Decode(&cycles)
	if err != nil {
		return nil, err
	}

	if len(cycles) == 0 {
		return nil, fmt.Errorf("app not found")
	}

	return cycles[0], nil
}

func (a *api) createApp(name string) (*data.App, error) {
	path, err := a.makePath("/apps")
	if err != nil {
		return nil, err
	}

	reqBody := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}

	resp, err := httpRequest("POST", path, reqBody, a.session)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, parseErrorMessage(resp.Body)
	}

	var app data.App

	err = json.NewDecoder(resp.Body).Decode(&app)
	if err != nil {
		return nil, err
	}

	return &app, nil
}

func (a *api) getAppByName(name string) (*data.App, error) {
	path, err := a.makePath("/apps?name=%s", appName)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest("GET", path, nil, a.session)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, parseErrorMessage(resp.Body)
	}

	var apps []*data.App

	err = json.NewDecoder(resp.Body).Decode(&apps)
	if err != nil {
		return nil, err
	}

	if len(apps) == 0 {
		return nil, fmt.Errorf("app not found")
	}

	return apps[0], nil
}

func (a *api) createCycle(appName, cycleName string) (*data.Cycle, error) {
	// first we need to get the app by name
	app, err := a.getAppByName(appName)
	if err != nil {
		return nil, err
	}

	// then we need to construct the api to point to that app and
	// create cycle

	path, err := a.makePath("/apps/%d/cycles", app.ID)
	if err != nil {
		return nil, err
	}

	reqBody := struct {
		Name string `json:"name"`
	}{
		Name: cycleName,
	}

	resp, err := httpRequest("POST", path, reqBody, a.session)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, parseErrorMessage(resp.Body)
	}

	var cycle *data.Cycle

	err = json.NewDecoder(resp.Body).Decode(&cycle)
	if err != nil {
		return nil, err
	}

	return cycle, nil
}

func (a *api) createRelease(appID, cycleID int64, platform, version, note string) (*data.Release, error) {
	path, err := a.makePath("apps/%d/cycles/%d/releases", appID, cycleID)
	if err != nil {
		return nil, err
	}

	reqBody := struct {
		Platform string `json:"platform"`
		Version  string `json:"version"`
		Note     string `json:"note"`
	}{
		Platform: platform,
		Version:  version,
		Note:     note,
	}

	resp, err := httpRequest("POST", path, reqBody, a.session)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, parseErrorMessage(resp.Body)
	}

	var release data.Release

	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return nil, err
	}

	return &release, nil
}

func (a *api) bundleUpload(appID, cycleID, releaseID int64, dataReader io.Reader) ([]*data.Bundle, error) {
	path, err := a.makePath("apps/%d/cycles/%d/releases/%d/bundles", appID, cycleID, releaseID)
	if err != nil {
		return nil, err
	}

	resp, err := httpRequest("POST", path, dataReader, a.session)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, parseErrorMessage(resp.Body)
	}

	var bundles []*data.Bundle

	err = json.NewDecoder(resp.Body).Decode(bundles)
	if err != nil {
		return nil, err
	}

	return bundles, nil
}

func newAPI(serverAddr string) *api {
	return &api{serverAddr: serverAddr}
}
