package handlers

import "net/http"

/* TODO: implement a CORS middleware handler, as described
in https://drstearns.github.io/tutorials/cors/ that responds
with the following headers to all requests:

  Access-Control-Allow-Origin: *
  Access-Control-Allow-Methods: GET, PUT, POST, PATCH, DELETE
  Access-Control-Allow-Headers: Content-Type, Authorization
  Access-Control-Expose-Headers: Authorization
  Access-Control-Max-Age: 600
*/

const methodsCORS = "Access-Control-Allow-Methods"
const allowHeadersCORS = "Access-Control-Allow-Headers"
const exposeHeadersCORS = "Access-Control-Expose-Headers"
const maxAgeCORS = "Access-Control-Max-Age"

func corsHandler(w http.ResponseWriter) {
	w.Header().Add(headerCORS, "*")
	w.Header().Add(methodsCORS, "GET, PUT, POST, PATCH, DELETE")
	w.Header().Add(allowHeadersCORS, "Content-Type, Authorization")
	w.Header().Add(exposeHeadersCORS, "Authorization")
	w.Header().Add(maxAgeCORS, "600")
}