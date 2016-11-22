package apps

import (
	"mime/multipart"
	"net/http"

	"path/filepath"

	"github.com/pressly/warpdrive"
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

				// save the file in temporary folder
				path := filepath.Join(warpdrive.Conf.Server.TempFolder, web.UUID())
				if err := web.CopyDataToFile(file, path); err != nil {
					return err
				}

				mapFiles[header.Filename] = path
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

}
