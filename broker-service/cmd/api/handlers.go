package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Emmanuel-MacAnThony/broker/events"
)

type RequestPayload struct {
	Action string `json:"action"`
	Auth AuthPayload `json:"auth,omitempty"`
	Log LogPayload `json:"log,omitempty"`
	Mail MailPayload `json:"mail,omitempty"`

}

type AuthPayload struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From string `json:"from"`
	To string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}


func (app *Config) Broker(w http.ResponseWriter, r *http.Request){

	payload := JSONResponse{
		Error: false,
		Message: "Hit thr broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

}

func (app *Config) HandleSubmission(w http.ResponseWriter,r *http.Request){
	      var requestPayload RequestPayload
		  err := app.readJSON(w,r,&requestPayload)
		  if err != nil {
			app.errorJSON(w, err, http.StatusBadRequest)
			return
		  }

		  switch requestPayload.Action {
		  case "auth":
			app.authenticate(w, requestPayload.Auth)
		  case "log":
			app.logEventViaRabbit(w, requestPayload.Log)
		  case "mail":
			app.sendMail(w, requestPayload.Mail)
		  default:
			app.errorJSON(w,errors.New("unknown action"))
		  }
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t") 
	request,err := http.NewRequest("POST","http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err :=client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	}else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("internal server error"))
		return
	}

	var jsonFromService JSONResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error{
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload JSONResponse
	payload.Error = false
	payload.Message = jsonFromService.Message
	payload.Data = jsonFromService.Data


	app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload){
	jsonData, _ :=json.MarshalIndent(entry, "", "\t")
	logServiceUrl := "http://logger-service/log"
	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated{
		app.errorJSON(w, err)
		return
	}

	var payload JSONResponse
	payload.Error = false
	payload.Message = "logged"

	app.writeJSON(w, http.StatusCreated, payload)

}


func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload){
	jsonData, _ := json.MarshalIndent(msg, "", "\t")  
	mailServiceURL := "http://mailer-service/send"
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil{
		app.errorJSON(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK{
		app.errorJSON(w, errors.New("error calling mail service"))
		return
	}

	var payload JSONResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload){

	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload JSONResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	app.writeJSON(w, http.StatusAccepted, payload)


}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := events.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j,_ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}

	return nil
}