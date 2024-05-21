package main

import (
	"log"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From string `json:"from"`
		To string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`

	}

	var requestPayload mailMessage
	err := app.readJSON(w,r,&requestPayload)
	if err!= nil{
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	msg := Message{
		From: requestPayload.From,
		To: requestPayload.To,
		Subject: requestPayload.Subject,
		Data: requestPayload.Message,
	}

	err = app.mailer.SendSMTPMessage(msg)
	if err!= nil{
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	payload := JSONResponse{
		Error: false,
		Message: "message sent to:" + msg.To,
	}

	app.writeJSON(w,http.StatusOK,payload)

}