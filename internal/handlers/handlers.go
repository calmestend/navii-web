package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/charmbracelet/log"
)

type DeletePageData struct {
	Errors []string
	Email  string
}

type SigninResponseData struct {
	AccessToken  string   `json:"access_token"`
	Error        string   `json:"error"`
	Messages[]     string `json:"message"`
}

type DeleteAccountResponseData struct {
	Message string `json:"message"`
	Status  int    `json:"statusCode"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ts, err := template.ParseFiles("./web/index.html")
	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func DeleteAccountForm(w http.ResponseWriter, r *http.Request) {
	err := renderDeleteAccountPage(w, DeletePageData{})
	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	confirmPassword := r.PostForm.Get("confirmPassword")

	if email == "" || password == "" || confirmPassword == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var errors []string

	if password != confirmPassword {
		errors = append(errors, "Las contraseñas no coinciden")
	}

	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)

		err := renderDeleteAccountPage(w, DeletePageData{Email: email, Errors: errors})
		if err != nil {
			log.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}

	url := "https://naviiapp.com/api/v1"

	payload := map[string]string{
		"email":    email,
		"password": password,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/local/signin", url), bytes.NewBuffer(body))
	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			log.Error("Closing response body", "err", err)
		}
	}(resp.Body)


	var signinResp SigninResponseData

	err = json.NewDecoder(resp.Body).Decode(&signinResp)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Error(signinResp.Messages)
		err := renderDeleteAccountPage(w, DeletePageData{Email: email, Errors: signinResp.Messages})
		if err != nil {
			log.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}

	log.Info("Signin Response Data", "Access Token", signinResp.AccessToken)

	req, err = http.NewRequest("POST", fmt.Sprintf("%s/auth/deleteAccount", url), nil)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+signinResp.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {

		var deleteAccountResp DeleteAccountResponseData

		err = json.NewDecoder(resp.Body).Decode(&deleteAccountResp)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		log.Error(deleteAccountResp.Message)
		err := renderDeleteAccountPage(w, DeletePageData{Email: email, Errors: []string{deleteAccountResp.Message}})
		if err != nil {
			log.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}
	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			log.Error("Closing response body", "err", err)
		}
	}(resp.Body)

	http.Redirect(w, r, "/", http.StatusFound)
}

func PrivacyPolicies(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./web/privacy_policies.html")
	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func renderDeleteAccountPage(w http.ResponseWriter, data DeletePageData) error {
	ts, err := template.ParseFiles("./web/user_delete.html")
	if err != nil {
		return err
	}

	err = ts.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}
