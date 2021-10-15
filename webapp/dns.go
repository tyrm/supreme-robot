package webapp

import (
	"github.com/tyrm/supreme-robot/models"
	"net/http"
)

type DnsPageTemplate struct {
	templateCommon

	Domains []models.Domain
}

func (s *Server) DnsGetHandler(w http.ResponseWriter, r *http.Request) {
	// Init template variables
	tmplVars := &DnsPageTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := r.Context().Value(UserKey).(*models.User)
	domains, err := s.db.ReadDomainsForUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmplVars.Domains = *domains

	err = s.templates.ExecuteTemplate(w, "dns", tmplVars)
	if err != nil {
		logger.Errorf("could not render dns template: %s", err.Error())
	}
}