package webapp

import (
	"github.com/tyrm/supreme-robot/config"
	"net/http"
)

type AdminDnsPageTemplate struct {
	templateCommon
	Breadcrumbs *[]templateBreadcrumb

	FormNS             *templateFormInput
	FormSectionGeneral *templateFormInput
	FormGeneralSubmit  *templateFormButton
}

func (s *Server) AdminDnsGetHandler(w http.ResponseWriter, r *http.Request) {
	// Init template variables
	tmplVars := &AdminDnsPageTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	tmplVars.PageTitle = "Users"
	tmplVars.Breadcrumbs = &[]templateBreadcrumb{
		{
			Text: "Admin Dashboard",
			HRef: "/app/admin",
		},
		{
			Text: "Users",
		},
	}

	tmplVars.FormGeneralSubmit = &templateFormButton{
		Color: "success",
		Text:  "Update",
	}

	// fetch config from
	confToFetch := []string{
		config.KeySoaNS,
	}
	confs, err := s.config.MGet(&confToFetch)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// general section
	tmplVars.FormNS = &templateFormInput{
		ID:          "primary-ns",
		Name:        "primary_ns",
		Placeholder: "ns1.example.com.",
		Required:    true,
	}
	if (*confs)[config.KeySoaNS] != nil {
		tmplVars.FormNS.Value = *(*confs)[config.KeySoaNS]
	}

	tmplVars.FormSectionGeneral = &templateFormInput{
		ID:    "section-general",
		Name:  "_section",
		Value: "general",
	}

	err = s.templates.ExecuteTemplate(w, "admin_dns", tmplVars)
	if err != nil {
		logger.Errorf("could not render dns template: %s", err.Error())
	}
}

func (s *Server) AdminDnsPostHandler(w http.ResponseWriter, r *http.Request) {
	// parse form data
	err := r.ParseForm()
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	switch r.Form.Get("_section") {
	case "general":
		confs := map[string]string{
			config.KeySoaNS: r.Form.Get("primary_ns"),
		}

		err = s.config.MSet(&confs)
		if err != nil {
			s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}

	default:
		s.returnErrorPage(w, r, http.StatusBadRequest, "missing section")
		return
	}

	// redirect to reload page
	http.Redirect(w, r, "/app/admin/dns", http.StatusFound)
}
