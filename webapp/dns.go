package webapp

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/tyrm/supreme-robot/config"
	"github.com/tyrm/supreme-robot/models"
	"net/http"
	"strconv"
)

type DnsPageTemplate struct {
	templateCommon

	Domains *[]models.Domain
}

type DnsDomainPageTemplate struct {
	templateCommon
	Breadcrumbs *[]templateBreadcrumb

	Domain  *models.Domain
	Records *[]models.Record
}

type DnsDomainFormTemplate struct {
	templateCommon
	Breadcrumbs *[]templateBreadcrumb

	TitleText      string
	FormId         *templateFormInput
	FormDomain     *templateFormInput
	FormSubmit     *templateFormButton
	FormSoaTTL     *templateFormInput
	FormSoaMBox    *templateFormInput
	FormSoaNS      *templateFormInput
	FormSoaRefresh *templateFormInput
	FormSoaRetry   *templateFormInput
	FormSoaExpire  *templateFormInput
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

	tmplVars.PageTitle = "Domains"

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
	s.displayDomainAddPage(w, r, "", "300", "", "604800", "86400", "2419200", "")
}

func (s *Server) DnsDomainAddPostHandler(w http.ResponseWriter, r *http.Request) {
	// parse form data
	err := r.ParseForm()
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// check if domain exists already
	dbDomain, err := s.db.ReadDomainByDomain(r.Form.Get("domain"))
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if dbDomain != nil {
		s.displayDomainAddPage(w, r, r.Form.Get("domain"), r.Form.Get("soa_ttl"), r.Form.Get("soa_mbox"), r.Form.Get("soa_refresh"), r.Form.Get("soa_retry"), r.Form.Get("soa_expire"), "domain exists")
		return
	}

	// TODO check if domain is parent of existing domain
	// TODO check if domain is child of existing domain

	// check domain validity
	user := r.Context().Value(UserKey).(*models.User)
	domain := models.Domain{
		Domain: r.Form.Get("domain"),
		OwnerID:  user.ID,
	}
	valid := domain.ValidateDomain()
	if !valid {
		s.displayDomainAddPage(w, r, r.Form.Get("domain"), r.Form.Get("soa_ttl"), r.Form.Get("soa_mbox"), r.Form.Get("soa_refresh"), r.Form.Get("soa_retry"), r.Form.Get("soa_expire"), "domain exists")
		return
	}

	// validation soa
	ttl, err := strconv.Atoi(r.Form.Get("soa_ttl"))
	if err != nil {
		s.displayDomainAddPage(w, r, r.Form.Get("domain"), r.Form.Get("soa_ttl"), r.Form.Get("soa_mbox"), r.Form.Get("soa_refresh"), r.Form.Get("soa_retry"), r.Form.Get("soa_expire"), "ttl is not integer")
		return
	}
	// TODO validate value of mbox
	refresh, err := strconv.Atoi(r.Form.Get("soa_refresh"))
	if err != nil {
		s.displayDomainAddPage(w, r, r.Form.Get("domain"), r.Form.Get("soa_ttl"), r.Form.Get("soa_mbox"), r.Form.Get("soa_refresh"), r.Form.Get("soa_retry"), r.Form.Get("soa_expire"), "refresh is not integer")
		return
	}
	retry, err := strconv.Atoi(r.Form.Get("soa_retry"))
	if err != nil {
		s.displayDomainAddPage(w, r, r.Form.Get("domain"), r.Form.Get("soa_ttl"), r.Form.Get("soa_mbox"), r.Form.Get("soa_refresh"), r.Form.Get("soa_retry"), r.Form.Get("soa_expire"), "retry is not integer")
		return
	}
	expire, err := strconv.Atoi(r.Form.Get("soa_expire"))
	if err != nil {
		s.displayDomainAddPage(w, r, r.Form.Get("domain"), r.Form.Get("soa_ttl"), r.Form.Get("soa_mbox"), r.Form.Get("soa_refresh"), r.Form.Get("soa_retry"), r.Form.Get("soa_expire"), "expire is not integer")
		return
	}

	// get server ns
	ns, err := s.config.Get(config.KeySoaNS)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if ns == nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, "NS not defined")
		return
	}

	// add to database
	err = domain.Create(s.db)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// create soa record
	record := models.Record{
		Name:     "@",
		DomainID: domain.ID,
		Type:     "SOA",
		Value:    *ns,
		TTL: sql.NullInt32{
			Int32: int32(ttl),
			Valid: true,
		},
		MBox: sql.NullString{
			String: r.Form.Get("soa_mbox"),
			Valid:  true,
		},
		Refresh: sql.NullInt32{
			Int32: int32(refresh),
			Valid: true,
		},
		Retry: sql.NullInt32{
			Int32: int32(retry),
			Valid: true,
		},
		Expire: sql.NullInt32{
			Int32: int32(expire),
			Valid: true,
		},
	}
	err = record.Create(s.db)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// schedule update
	err = s.scheduler.AddDomain(domain.ID)
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

func (s *Server) DnsDomainDeleteGetHandler(w http.ResponseWriter, r *http.Request) {
	// get requested domain
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		s.returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}
	domain, err := s.db.ReadDomain(id)
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

	s.displayDomainDeletePage(w, r, domain.Domain, "")
}

func (s *Server) DnsDomainDeletePostHandler(w http.ResponseWriter, r *http.Request) {
	// get requested domain
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		s.returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}
	domain, err := s.db.ReadDomain(id)
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

	// do delete
	err = domain.Delete(s.db)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// schedule update
	err = s.scheduler.RemoveDomain(domain.ID)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// redirect to domain page
	us := r.Context().Value(SessionKey).(*sessions.Session)
	us.Values["page-alert-success"] = templateAlert{Text: "Domain deleted"}
	err = us.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/app/dns", http.StatusFound)
}

func (s *Server) DnsDomainGetHandler(w http.ResponseWriter, r *http.Request) {
	// get requested domain
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		s.returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}
	domain, err := s.db.ReadDomain(id)
	if err != nil {
		logger.Errorf("db error: %s", err.Error())
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
	records, err := domain.GetRecords(s.db)
	if err != nil {
		logger.Errorf("db error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmplVars.Records = records

	err = s.templates.ExecuteTemplate(w, "dns_domain", tmplVars)
	if err != nil {
		logger.Errorf("could not render template: %s", err.Error())
	}
}

func (s *Server) displayDomainAddPage(w http.ResponseWriter, r *http.Request, domain, soaTtl, soaMBox, soaRefresh, soaRetry, soaExpire, formError string) {
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

	tmplVars.FormSoaTTL = &templateFormInput{
		ID:          "soa_ttl",
		Name:        "soa_ttl",
		Value:       soaTtl,
		Placeholder: "300",
		Required:    true,
	}
	tmplVars.FormSoaMBox = &templateFormInput{
		ID:          "soa_mbox",
		Name:        "soa_mbox",
		Value:       soaMBox,
		Placeholder: "hostmaster.example.com.",
		Required:    true,
	}
	tmplVars.FormSoaRefresh = &templateFormInput{
		ID:          "soa_refresh",
		Name:        "soa_refresh",
		Value:       soaRefresh,
		Placeholder: "604800",
		Required:    true,
	}
	tmplVars.FormSoaRetry = &templateFormInput{
		ID:          "soa_retry",
		Name:        "soa_retry",
		Value:       soaRetry,
		Placeholder: "86400",
		Required:    true,
	}
	tmplVars.FormSoaExpire = &templateFormInput{
		ID:          "soa_expire",
		Name:        "soa_expire",
		Value:       soaExpire,
		Placeholder: "2419200",
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

func (s *Server) displayDomainDeletePage(w http.ResponseWriter, r *http.Request, domain, formError string) {
	// Init template variables
	tmplVars := &DnsDomainFormTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmplVars.PageTitle = "Delete Domain"
	tmplVars.Breadcrumbs = &[]templateBreadcrumb{
		{
			Text: "DNS Manager",
			HRef: "/app/dns",
		},
		{
			Text: "Delete Domain",
		},
	}

	tmplVars.TitleText = "Edit Domain"
	tmplVars.FormDomain = &templateFormInput{
		ID:          "domain",
		Name:        "domain",
		Value:       domain,
		Placeholder: "example.com.",
		Disabled:    true,
	}
	tmplVars.FormSubmit = &templateFormButton{
		Color: "danger",
		Text:  "Delete",
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
