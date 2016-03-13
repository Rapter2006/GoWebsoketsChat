package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type JsonErr struct {
	Code int    `json:"code"`
	Text string `json:"desc"`
}

func NewErrorRender(code int, text string, w http.ResponseWriter) {
	errorMessage := &JsonErr{
		Code: code,
		Text: text,
	}
	errorMessage.NewErrorRender(w)
}

func (e *JsonErr) NewErrorRender(w http.ResponseWriter) {
	fmt.Println("Error!")
	fmt.Fprintf(
		logFile,
		"%s\t%s\n",
		strconv.Itoa(e.Code),
		e.Text)

	json.NewEncoder(w).Encode(e)
}
