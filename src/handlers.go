package main

import (
	"net/http"
	"github.com/gorilla/mux"
)

func CreateWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomId := vars["roomId"]
	userId := vars["userId"]

	if roomId == "" || userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		NewErrorRender(1, "Параметр не указан", w)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		NewErrorRender(2, "Вебсокет соединение не работает", w)
		return
	}

	ConnectionInitialize(ws, &roomId, &userId)
}

