package apps

import (
	"mime/multipart"
	"net/http"

	"path/filepath"

	"github.com/pressly/warpdrive"
	"github.com/pressly/warpdrive/lib/folder"
	"github.com/pressly/warpdrive/lib/warp"
	"github.com/pressly/warpdrive/services"
	"github.com/pressly/warpdrive/web"
)

func saveFilesAsTemporary(files map[string][]*multipart.FileHeader) (map[string]string, error) {
	//mapFiles is a map which key represents the real filename and the value
	//represents the temporary location of the file
	mapFiles := make(map[string]string)

	for _, headers := range files {
		for _, header := range headers {
			//we create a clouser here so we can use defer to close the file
			//in case of any errors.
			err := func(header *multipart.FileHeader) error {
				file, _ := header.Open()
				defer file.Close()

				// create a temporary folder
				path := filepath.Join(warpdrive.Conf.Server.TempFolder, web.UUID())
				// extract tar.gz file to the newly created teamporary file
				if err := warp.Extract(file, path); err != nil {
					return err
				}

				// grab all the files under the newly created temporary folder
				// which contains the extracted received tar.gz file
				filePaths, err := folder.ListFilePaths(path)
				if err != nil {
					return err
				}

				// need to build a map <filename> -> file path
				for _, filePath := range filePaths {
					mapFiles[filepath.Base(filePath)] = filepath.Join(path, filePath)
				}

				return nil
			}(header)

			if err != nil {
				return nil, err
			}
		}
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

	err = r.ParseMultipartForm(warpdrive.Conf.FileUpload.FileMaxSize)
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	mapFiles, err := saveFilesAsTemporary(r.MultipartForm.File)
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	bundles, err := services.CreateBundles(userID, appID, cycleID, releaseID, mapFiles)
	if err != nil {
		web.Respond(w, http.StatusOK, err)
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
