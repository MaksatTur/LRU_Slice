package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"lrumemcache/data"
	"net/http"
)

//Service ...
type Service struct {
	ServerURL string
	LRU       *data.LRU
}

//NewService ...
func NewService(serverURL string, capacityLRU int) *Service {
	return &Service{
		ServerURL: serverURL,
		LRU:       data.NewLRU(capacityLRU),
	}
}

func (s *Service) mainHandler(w http.ResponseWriter, r *http.Request) {
	fileContents, err := ioutil.ReadFile("./static/index.html")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write(fileContents)
}

func (s Service) todosHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	switch r.Method {
	case http.MethodGet:
		fmt.Printf("From Get todos = %v\n", s.LRU.ToDoItems)
		todosJSON, err := json.Marshal(s.LRU.ToDoItems)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(todosJSON)
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		todo := data.ToDo{}
		err := decoder.Decode(&todo)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		s.LRU.Set(todo.Name, todo)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Service) panicMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recovered", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

//StartService ...
func (s *Service) StartService() {
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/", s.mainHandler)
	apiMux.HandleFunc("/todos/", s.todosHandler)
	mainHandler := s.panicMiddlware(apiMux)
	fmt.Printf("Service is started on %s\n", s.ServerURL)
	log.Fatal(http.ListenAndServe(s.ServerURL, mainHandler))
}
