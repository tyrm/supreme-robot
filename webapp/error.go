package webapp

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorPageTemplate struct {
	templateCommon

	BotImage    string
	Header      string
	SubHeader   string
	Paragraph   string
	ButtonHRef  string
	ButtonLabel string
}

var (
	ErrInvalidCount     = errors.New("value of count must be greater than zero")
	ErrInvalidPage      = errors.New("value of page must be greater than zero")
	ErrUnknownAttribute = errors.New("unknown attribute")
)

func (s *Server) returnErrorPage(w http.ResponseWriter, r *http.Request, code int, errStr string) {

	// Init template variables
	tmplVars := &ErrorPageTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// change CSS sheet
	tmplVars.HeadCSS = &[]templateHeadLink{
		{
			HRef: "/static/css/error.css",
			Rel:  "stylesheet",
		},
	}

	// custom body css
	tmplVars.BodyClass = "text-center"

	// set bot image
	switch code {
	case http.StatusBadRequest:
		// 400
		tmplVars.BotImage = BotEmojiConfused
	case http.StatusUnauthorized:
		// 401
		tmplVars.BotImage = BotEmojiAngry
	case http.StatusForbidden:
		// 403
		tmplVars.BotImage = BotEmojiMad
	case http.StatusNotFound:
		// 404
		tmplVars.BotImage = BotEmojiConfused
	case http.StatusMethodNotAllowed:
		// 405
		tmplVars.BotImage = BotEmojiMad
	default:
		tmplVars.BotImage = BotEmojiOffline
	}

	// set text
	tmplVars.Header = fmt.Sprintf("%d", code)
	tmplVars.SubHeader = http.StatusText(code)
	tmplVars.Paragraph = errStr

	// set top button
	switch code {
	case http.StatusUnauthorized:
		tmplVars.ButtonHRef = "/login"
		tmplVars.ButtonLabel = "Login"
	default:
		tmplVars.ButtonHRef = "/app/"
		tmplVars.ButtonLabel = "Home"
	}

	w.WriteHeader(code)
	err = s.templates.ExecuteTemplate(w, "error", tmplVars)
	if err != nil {
		logger.Errorf("could not render home template: %s", err.Error())
	}
}

func (s *Server) MethodNotAllowedHandler() http.Handler {
	// wrap in middleware since middlware isn't run on error pages
	return s.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.returnErrorPage(w, r, http.StatusMethodNotAllowed, "")
	}))
}

func (s *Server) NotFoundHandler() http.Handler {
	// wrap in middleware since middlware isn't run on error pages
	return s.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("page not found: %s", r.URL.Path))
	}))
}
