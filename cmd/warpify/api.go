package warpify

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/pressly/warpdrive/config"
	"github.com/pressly/warpdrive/data"
	"github.com/pressly/warpdrive/lib/crypto"
	"github.com/pressly/warpdrive/lib/httpclient"
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
	warpFile *config.ClientConfig
}

func (a *api) makePath(path string, args ...interface{}) (string, error) {
	path, err := httpclient.JoinURL(a.warpFile.ServerAddr, fmt.Sprintf(path, args...))
	if err != nil {
		return "", fmt.Errorf("Server Address '%s' is invalid", a.warpFile.ServerAddr)
	}
	return path, nil
}

func (a *api) checkVersion(appID, cycleID int64, platform, currentVersion string) (map[string]*data.Release, error) {
	path, err := a.makePath("/apps/%d/cycles/%d/version/%s/platform/%s/latest", appID, cycleID, currentVersion, platform)
	if err != nil {
		return nil, err
	}

	resp, err := httpclient.Request("GET", path, nil, "", "")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, parseErrorMessage(resp.Body)
	}

	// releaseMap has 2 keys, soft and hard.
	// soft means, you can update it, hard means update available via app store or play store
	releaseMap := make(map[string]*data.Release)

	err = json.NewDecoder(resp.Body).Decode(&releaseMap)
	if err != nil {
		return nil, err
	}

	return releaseMap, nil
}

func (a *api) remoteVersions(appID, cycleID int64) ([]*data.Release, error) {
	// prepare the path
	path, err := a.makePath("/apps/%d/cycles/%d/releases?platform=%s", appID, cycleID, conf.platform)
	if err != nil {
		return nil, err
	}

	resp, err := httpclient.Request("GET", path, nil, "", "")
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

func (a *api) downloadRelease(appID, cycleID, releaseID int64) (io.ReadCloser, error) {
	// check if appID and cycleID exists in warpFile and also
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

	resp, err := httpclient.Request("POST", path, encryptedKey, "", "")
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
		err := crypto.AESDecryptStream(key, resp.Body, w)
		if err != nil {
			w.CloseWithError(err)
		} else {
			w.Close()
		}
	}()

	return r, nil
}

func (a *api) downloadVersion(appID, cycleID int64, version, platform string) (io.ReadCloser, error) {
	// check if appID and cycleID exists in warpFile and also
	// we need cycle to grab ke public key so we can encrypt the session key
	cycleConfig, err := a.warpFile.GetCycleByID(appID, cycleID)
	if err != nil {
		return nil, err
	}

	// prepare the path
	path, err := a.makePath("/apps/%d/cycles/%d/version/%s/platform/%s/download", appID, cycleID, version, platform)
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

	resp, err := httpclient.Request("POST", path, encryptedKey, "", "")
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
		err := crypto.AESDecryptStream(key, resp.Body, w)
		if err != nil {
			w.CloseWithError(err)
		} else {
			w.Close()
		}
	}()

	return r, nil
}

func makeAPI(warpFile *config.ClientConfig) *api {
	return &api{
		warpFile: warpFile,
	}
}
