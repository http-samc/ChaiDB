package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type server struct{}

type target struct {
	endpoint string
	structure []string
}

func (t *target) GetStructure() {
	t.structure = strings.Split(t.endpoint,"/")[1:]
}

func (t *target) QueryingFavicon() bool {
	if t.structure[0] == "favicon.ico" { return true }
	return false
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	t := &target{r.RequestURI, []string{}}
	t.GetStructure()
	if t.QueryingFavicon() {
		w.WriteHeader(http.StatusOK)
	} else if len(t.structure) == 1 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Welcome to the API Base"}`))
	} else if len(t.structure) > 0 {
		fmt.Println(t.structure)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"!!!"}`))
	}
}

func main() {
	s := &server{}
	http.Handle("/", s)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}