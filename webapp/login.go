package webapp

import (
	"github.com/gorilla/sessions"
	"net/http"
)

type LoginPageTemplate struct {
	templateCommon

	BotImage string

	FormError    string
	FormUsername string
	FormPassword string
}

func (s *Server) LoginGetHandler(w http.ResponseWriter, r *http.Request) {
	s.displayLoginPage(w, r, "", "", BotEmojiHappy, "")
}

func (s *Server) LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	// parse form data
	err := r.ParseForm()
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// check if user exists
	formUsername := r.Form.Get("username")
	formPassword := r.Form.Get("password")
	user, err := s.db.ReadUserByUsername(formUsername)
	if err != nil {
		logger.Errorf("db error: %s", err.Error())
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if user == nil {
		s.displayLoginPage(w, r, formUsername, formPassword, BotEmojiMad, LoginErrorTest)
		return
	}

	// check password validity
	passValid := user.CheckPasswordHash(formPassword)
	if passValid == false {
		s.displayLoginPage(w, r, formUsername, formPassword, BotEmojiMad, LoginErrorTest)
		return
	}

	// Init Session
	us := r.Context().Value(SessionKey).(*sessions.Session)
	us.Values["user"] = user
	err = us.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect to last page
	val := us.Values["login-redirect"]
	var loginRedirect string
	var ok bool
	if loginRedirect, ok = val.(string); !ok {
		// redirect home page if no login-redirect
		http.Redirect(w, r, "/app/", http.StatusFound)
		return
	}

	http.Redirect(w, r, loginRedirect, http.StatusFound)
	return
}

func (s *Server) displayLoginPage(w http.ResponseWriter, r *http.Request, username, password, botImage, formError string) {
	// Init template variables
	tmplVars := &LoginPageTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// change CSS sheet
	tmplVars.HeadCSS = &[]templateHeadLink{
		{
			HRef: "/static/css/login.css",
			Rel:  "stylesheet",
		},
	}

	tmplVars.PageTitle = "Login"

	// set bot image
	tmplVars.BotImage = botImage

	// set form values
	tmplVars.FormError = formError
	tmplVars.FormUsername = username
	tmplVars.FormPassword = password

	// custom body css
	tmplVars.BodyClass = "text-center"

	err = s.templates.ExecuteTemplate(w, "login", tmplVars)
	if err != nil {
		logger.Errorf("could not render home template: %s", err.Error())
	}
}
