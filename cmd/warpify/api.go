package warpify

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/pressly/warpdrive/config"
	"github.com/pressly/warpdrive/data"
	"github.com/pressly/warpdrive/lib/crypto"
	"github.com/pressly/warpdrive/lib/rest"
)

func parseErrorMessage(body io.Reader) error {
	var err error

	errorMessage := struct {
		Error string `json:"error"`
	}{}

	err = json.NewDecoder(body).Decode(&errorMessage)
	if err != nil {
		errorMessage.Error = err.Error()
	}

	return fmt.Errorf(errorMessage.Error)
}

type api struct {
	warpFile   *config.ClientConfig
	versionMap *config.VersionMap
}

func (a *api) makePath(path string, args ...interface{}) (string, error) {
	path, err := rest.JoinURL(a.warpFile.ServerAddr, fmt.Sprintf(path, args...))
	if err != nil {
		return "", fmt.Errorf("Server Address '%s' is invalid", a.warpFile.ServerAddr)
	}
	return path, nil
}

func (a *api) checkVersion(appID, cycleID int64, platform, currentVersion string) ([]*data.Release, error) {
	path, err := a.makePath("/apps/%d/cycles/%d/releases/latest/version/%s/platform/%s", appID, cycleID, currentVersion, platform)
	if err != nil {
		return nil, err
	}

	resp, err := rest.Request("GET", path, nil, "", "")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, parseErrorMessage(resp.Body)
	}

	var releases []*data.Release

	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		return nil, err
	}

	return releases, nil
}

func (a *api) downloadVersion(appID, cycleID, releaseID int64) (io.Reader, error) {
	// check if appID and cycleID exists in warpFile
	// we need cycle to grab ke public key so we can encrypt the session key
	cycleConfig, err := a.warpFile.GetCycleByID(appID, cycleID)
	if err != nil {
		return nil, err
	}

	// prepare the path
	path, err := a.makePath("/apps/%d/cycles/%d/releases/%d/download", appID, cycleID, releaseID)
	if err != nil {
		return nil, err
	}

	// we need to create a session key here
	key, err := crypto.MakeAESKey(32)
	if err != nil {
		return nil, err
	}

	// if we have this far it means that everything is correct and
	// we need to encrypt the session key with public key
	encryptedKey, err := crypto.EncryptByPublicRSA(key, cycleConfig.Key, "sha256")
	if err != nil {
		return nil, err
	}

	resp, err := rest.Request("POST", path, encryptedKey, "", "")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, parseErrorMessage(resp.Body)
	}

	// by now we have the stream of encrypted data which we need to decrypt by
	// session key that we have
	r, w := io.Pipe()
	go func() {
		crypto.AESDecryptStream(key, resp.Body, w)
	}()

	return r, nil
}

func makeApi(warpFile *config.ClientConfig, versionMap *config.VersionMap) *api {
	return &api{
		warpFile:   warpFile,
		versionMap: versionMap,
	}
}
