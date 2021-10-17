package webapp

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/tyrm/supreme-robot/models"
	"net/http"
)

type DnsPageTemplate struct {
	templateCommon

	Domains *[]models.Domain
}


type DnsDomainPageTemplate struct {
	templateCommon
	Breadcrumbs *[]templateBreadcrumb

	Domain *models.Domain
}

type DnsDomainFormTemplate struct {
	templateCommon
	Breadcrumbs *[]templateBreadcrumb

	TitleText  string
	FormId     *templateFormInput
	FormDomain *templateFormInput
	FormSubmit *templateFormButton
}

var (
	dnsDefaultOrder = "domain"
	dnsOrderMap     = map[string]string{
		"created_at": "created_at",
		"domain":     "domain",
	}
)

func (s *Server) DnsGetHandler(w http.ResponseWriter, r *http.Request) {
	// Init template variables
	tmplVars := &DnsPageTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// handle pagination
	page, count, _, err := paginationFromRequest(r, 50)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// get order by value
	orderBy, err := orderByFromRequest(r, dnsOrderMap, dnsDefaultOrder)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// get order direction
	orderAsc, err := orderAscFromRequest(r)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	user := r.Context().Value(UserKey).(*models.User)
	domains, err := s.db.ReadDomainsPageForUser(user, int(page-1), int(count), orderBy, orderAsc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmplVars.Domains = domains

	err = s.templates.ExecuteTemplate(w, "dns", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}

func (s *Server) DnsDomainAddGetHandler(w http.ResponseWriter, r *http.Request) {
	s.displayAddPage(w, r, "", "")
}

func (s *Server) DnsDomainAddPostHandler(w http.ResponseWriter, r *http.Request) {
	// parse form data
	err := r.ParseForm()
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// check domain validity
	user := r.Context().Value(UserKey).(*models.User)
	domain := models.Domain{
		Domain: r.Form.Get("domain"),
		Owner:  user,
	}
	valid := domain.ValidateDomain()
	if !valid {
		s.displayAddPage(w, r, r.Form.Get("domain"), "invalid domain")
		return
	}

	// add to database
	err = domain.Create(s.db)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// redirect to domain page
	us := r.Context().Value(SessionKey).(*sessions.Session)
	us.Values["page-alert-success"] = templateAlert{Text: "Domain added"}
	err = us.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect home
	http.Redirect(w, r, fmt.Sprintf("/app/dns/%s", domain.ID), http.StatusFound)
}

func (s *Server) DnsDomainGetHandler(w http.ResponseWriter, r *http.Request) {
	// get requested domain
	vars := mux.Vars(r)
	domain, err := s.db.ReadDomain(vars["id"])
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if domain == nil {
		s.returnErrorPage(w, r, http.StatusNotFound, "domain not found")
		return
	}

	// does the user own this domain
	user := r.Context().Value(UserKey).(*models.User)
	if domain.OwnerID != user.ID {
		s.returnErrorPage(w, r, http.StatusUnauthorized, "you don't own that domain")
		return
	}

	// Init template variables
	tmplVars := &DnsDomainPageTemplate{}
	err = initTemplate(w, r, tmplVars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmplVars.PageTitle = fmt.Sprintf("Domain %s", domain.Domain)
	tmplVars.Breadcrumbs = &[]templateBreadcrumb{
		{
			Text: "DNS Manager",
			HRef: "/app/dns",
		},
		{
			Text: fmt.Sprintf("Domain %s", domain.Domain),
		},
	}

	tmplVars.Domain = domain

	err = s.templates.ExecuteTemplate(w, "dns_domain", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}

func (s *Server) displayAddPage(w http.ResponseWriter, r *http.Request, domain, formError string) {
	// Init template variables
	tmplVars := &DnsDomainFormTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmplVars.PageTitle = "Add Domain"
	tmplVars.Breadcrumbs = &[]templateBreadcrumb{
		{
			Text: "DNS Manager",
			HRef: "/app/dns",
		},
		{
			Text: "Add Domain",
		},
	}

	tmplVars.TitleText = "Add Domain"
	tmplVars.FormDomain = &templateFormInput{
		ID:          "domain",
		Name:        "domain",
		Value:       domain,
		Placeholder: "example.com.",
		Required:    true,
	}
	tmplVars.FormSubmit = &templateFormButton{
		Color: "success",
		Text:  "Add",
	}

	if formError != "" {
		tmplVars.AlertError = &templateAlert{
			Text: formError,
		}
	}

	err = s.templates.ExecuteTemplate(w, "dns_domain_form", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}

func (s *Server) displayEditPage(w http.ResponseWriter, r *http.Request, domain, formError string) {
	// Init template variables
	tmplVars := &DnsDomainFormTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmplVars.PageTitle = "Edit Domain"
	tmplVars.Breadcrumbs = &[]templateBreadcrumb{
		{
			Text: "DNS Manager",
			HRef: "/dns",
		},
		{
			Text: "Add Domain",
		},
	}

	tmplVars.TitleText = "Edit Domain"
	tmplVars.FormDomain = &templateFormInput{
		ID:          "domain",
		Name:        "domain",
		Value:       domain,
		Placeholder: "example.com.",
		Required:    true,
	}
	tmplVars.FormSubmit = &templateFormButton{
		Color: "success",
		Text:  "Add",
	}

	err = s.templates.ExecuteTemplate(w, "dns_domain_form", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}
