package api

import (
	"fmt"
	"lrumemcache/data"
	"lrumemcache/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTodos(t *testing.T) {
	todo := data.ToDo{Name: "test", Done: false}
	config := utils.NewConfig()
	service := NewService(config.ServerURL, config.Capacity)
	service.LRU.Set(todo.Name, todo)
	handler := http.HandlerFunc(service.TodosHandler)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s%s", "http://", config.ServerURL, "/todos/"), nil)
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Wrong code. Expected %d, got %d", http.StatusOK, w.Code)
	}

	expected := `{"test":{"name":"test","done":false}}`
	if expected != w.Body.String() {
		t.Errorf("expected %s, got %s", expected, w.Body.String())
	}
}
