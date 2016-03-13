package main

import "net/http"

// Структурка, которая отвечает за описание одного роута
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"CreateWebsocketConnection",
		"GET",
		"/chat/{roomId}/{userId}",
		CreateWebsocketConnection,
	},
}
