package apps

import (
	"io"
	"net/http"

	"path/filepath"

	"strings"

	"github.com/pressly/warpdrive"
	"github.com/pressly/warpdrive/lib/folder"
	"github.com/pressly/warpdrive/lib/warp"
	"github.com/pressly/warpdrive/services"
	"github.com/pressly/warpdrive/web"
)

func saveFilesAsTemporary(reader io.ReadCloser) (map[string]string, error) {
	// create a temporary folder
	path := filepath.Join(warpdrive.Conf.Server.TempFolder, web.UUID())
	if err := warp.Extract(reader, path); err != nil {
		return nil, err
	}

	// grab all the files under the newly created temporary folder
	// which contains the extracted received tar.gz file
	filePaths, err := folder.ListFilePaths(path)
	if err != nil {
		return nil, err
	}

	//mapFiles is a map which key represents the real filename and the value
	//represents the temporary location of the file
	mapFiles := make(map[string]string)

	// we need to remove the temp folder from file path
	// alos path is temporary folder but we need to add slash at the end of it
	// so when we remove the path, we remove the / as well
	path = path + "/"

	// need to build a map <filename> -> file path
	for _, filePath := range filePaths {
		mapFiles[strings.Replace(filePath, path, "", 1)] = filePath
	}

	return mapFiles, nil
}

func uploadBundlesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user:id").(int64)

	appID, err := web.ParamAsInt64(r, "appId")
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	cycleID, err := web.ParamAsInt64(r, "cycleId")
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	releaseID, err := web.ParamAsInt64(r, "releaseId")
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	mapFiles, err := saveFilesAsTemporary(r.Body)
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	bundles, err := services.CreateBundles(userID, appID, cycleID, releaseID, mapFiles)
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	web.Respond(w, http.StatusOK, bundles)
}

func getBundlesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user:id").(int64)

	appID, err := web.ParamAsInt64(r, "appId")
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	cycleID, err := web.ParamAsInt64(r, "cycleId")
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	releaseID, err := web.ParamAsInt64(r, "releaseId")
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	query := r.URL.Query()
	name := query.Get("name")

	bundles, err := services.SearchBundles(userID, appID, cycleID, releaseID, name)
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	web.Respond(w, http.StatusOK, bundles)
}
