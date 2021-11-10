package graphql

import (
	"net/http"
	"net/url"
	"testing"
)

func TestMethodNotAllowedHandler(t *testing.T) {
	// create server
	server, _, _, _, err := newTestServer()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}
	if server == nil {
		t.Errorf("expected server, got: nil, want: *Server.")
	}

	handler := server.methodNotAllowedHandler()

	req := http.Request{
		Method: "POST",
		Proto:  "HTTP/2",
		URL: &url.URL{
			Path: "/",
		},
	}
	res := mockResponseWriter{
		Headr: http.Header{},
	}
	handler.ServeHTTP(&res, &req)

	expectedResponse := "{\"error\":{\"detail\":\"POST\",\"status\":405,\"title\":\"Method Not Allowed\"}}\n"
	if string(res.Body) != expectedResponse {
		t.Errorf("invalid body, got: %s, want: %s.", string(res.Body), expectedResponse)
	}
}

func TestNotFoundHandler(t *testing.T) {
	// create server
	server, _, _, _, err := newTestServer()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}
	if server == nil {
		t.Errorf("expected server, got: nil, want: *Server.")
	}

	handler := server.notFoundHandler()

	req := http.Request{
		Method: "POST",
		Proto:  "HTTP/2",
		URL: &url.URL{
			Path: "/home",
		},
	}
	res := mockResponseWriter{
		Headr: http.Header{},
	}
	handler.ServeHTTP(&res, &req)

	expectedResponse := "{\"error\":{\"detail\":\"page not found: /home\",\"status\":404,\"title\":\"Not Found\"}}\n"
	if string(res.Body) != expectedResponse {
		t.Errorf("invalid body, got: %s, want: %s.", string(res.Body), expectedResponse)
	}
}

func TestReturnErrorPage(t *testing.T) {
	// create server
	server, _, _, _, err := newTestServer()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}
	if server == nil {
		t.Errorf("expected server, got: nil, want: *Server.")
	}

	res := mockResponseWriter{}
	server.returnErrorPage(&res, 500, "error text")

	expectedResponse := "{\"error\":{\"detail\":\"error text\",\"status\":500,\"title\":\"Internal Server Error\"}}\n"
	if string(res.Body) != expectedResponse {
		t.Errorf("invalid body, got: %s, want: %s.", string(res.Body), expectedResponse)
	}
}
