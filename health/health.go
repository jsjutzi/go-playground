package health

import (
	"fmt"
	"net/http"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is up and running!")
}
