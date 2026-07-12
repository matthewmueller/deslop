package example

import "net/http"

func Bad() {
	http.NewServeMux() // want "do not use http.NewServeMux"
}
