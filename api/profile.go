package api

import (
	"io/ioutil"
	"net/http"
)

func (api *API) ImgProfileView(w http.ResponseWriter, r *http.Request) {
	// View with response image `img-avatar.png` from path `assets/images`
	// TODO: answer here
	read, _ := ioutil.ReadFile("assets/images/img-avatar.png")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(read)
	// w.Write(fileBytes) // fileBytes is []byte

}

func (api *API) ImgProfileUpdate(w http.ResponseWriter, r *http.Request) {
	// Update image `img-avatar.png` from path `assets/images`
	// TODO: answer here
	file, _, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		data, err := ioutil.ReadAll(file)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = ioutil.WriteFile("assets/images/img-avatar.png", data, 0777)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	api.dashboardView(w, r)
}
