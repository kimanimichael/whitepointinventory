package main

import "net/http"

func handlerStaticFiles(w http.ResponseWriter, r *http.Request) {

	staticDir := "./assets"

	
	http.ServeFile(w, r, staticDir+"/index.html")

	respondWithoutJSON(w, 200, struct{}{})
}