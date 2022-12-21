package api

import (
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"net/http"
	"path"
	"text/template"
	"time"

	"github.com/google/uuid"
)

func (api *API) Register(w http.ResponseWriter, r *http.Request) {
	// Read username and password request with FormValue.
	// creds := model.Credentials{} // TODO: replace this

	// if r.Method != http.MethodPost {
	// 	w.WriteHeader(405)
	// 	var error model.ErrorResponse
	// 	error.Error = "Method is not allowed!"
	// 	json.NewEncoder(w).Encode(error)
	// 	// kembali ke page register

	// } else {
	creds := model.Credentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	if len(creds.Password) == 0 || len(creds.Username) == 0 {
		// SendErrorJson(w, "Username or Password empty", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		var error model.ErrorResponse
		error.Error = "Username or Password empty"
		json.NewEncoder(w).Encode(error)
		return
	}

	// }
	// Handle request if creds is empty send response code 400, and message "Username or Password empty"
	// TODO: answer here

	err := api.usersRepo.AddUser(creds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	filepath := path.Join("views", "status.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	var data = map[string]string{"name": creds.Username, "message": "register success!"}
	err = tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
	}
}

func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	var resError model.ErrorResponse
	// Read usernmae and password request with FormValue.
	creds := model.Credentials{} // TODO: replace this

	creds.Username = r.FormValue("username")
	creds.Password = r.FormValue("password")

	// Handle request if creds is empty send response code 400, and message "Username or Password empty"
	// TODO: answer here
	// TODO: answer here
	if creds.Username == "" || creds.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")

		resError.Error = "Username or Password empty"
		js, _ := json.Marshal(resError)
		w.Write(js)
		return
	}
	err := api.usersRepo.LoginValid(creds)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Wrong User or Password!"})

		return
	}
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(5 * time.Hour)
	// Generate Cookie with Name "session_token", Path "/", Value "uuid generated with github.com/google/uuid", Expires time to 5 Hour.
	// TODO: answer here
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Path:    "/",
		Expires: expiresAt, // kita juga menetapkan expired 120 detik
	})

	session := model.Session{} // TODO: replace this
	session.Token = sessionToken
	session.Expiry = expiresAt
	session.Username = creds.Username

	err = api.sessionsRepo.AddSessions(session)
	if err != nil {
		panic(err)
	}

	api.dashboardView(w, r)
}

func SendErrorJson(w http.ResponseWriter, s string, i int) {
	panic("unimplemented")
}
func (api *API) Logout(w http.ResponseWriter, r *http.Request) {
	//Read session_token and get Value:
	// sessionToken := "" // TODO: replace this
	sessionToken, _ := r.Cookie("session_token")

	// Delete session from database
	api.sessionsRepo.DeleteSessions(sessionToken.Value)
	// if err != nil {
	// 	// w.WriteHeader(http.StatusInternalServerError)
	// 	// json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
	// 	return
	// }
	// // api.sessionsRepo.DeleteSessions(sessionToken)

	//Set Cookie name session_token value to empty and set expires time to Now:
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
		Path:    "/",
	})

	// TODO: answer here

	filepath := path.Join("views", "login.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
	}
}
