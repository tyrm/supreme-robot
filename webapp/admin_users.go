package webapp

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/tyrm/supreme-robot/models"
	"net/http"
)

type AdminUserPageTemplate struct {
	templateCommon
	Breadcrumbs *[]templateBreadcrumb

	User *models.User
}

type AdminUsersPageTemplate struct {
	templateCommon
	Breadcrumbs *[]templateBreadcrumb

	Users []models.User
}

type AdminUsersFormTemplate struct {
	templateCommon
	Breadcrumbs *[]templateBreadcrumb

	TitleText    string
	FormId       *templateFormInput
	FormUsername *templateFormInput
	FormPassword *templateFormInput
	FormSubmit   *templateFormButton
}

var (
	userDefaultOrder = "username"
	userOrderMap     = map[string]string{
		"created_at": "created_at",
		"username":   "username",
	}
)

func (s *Server) AdminUsersGetHandler(w http.ResponseWriter, r *http.Request) {
	if u := r.Context().Value(UserKey).(*models.User); !u.IsMemberOfGroup(&models.GroupsUserAdmin) {
		s.returnErrorPage(w, r, http.StatusUnauthorized, "You aren't authorized")
		return
	}

	// Init template variables
	tmplVars := &AdminUsersPageTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	// handle pagination
	page, count, _, err := paginationFromRequest(r, 50)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// get order by value
	orderBy, err := orderByFromRequest(r, userOrderMap, userDefaultOrder)
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

	userList, err := s.db.ReadUsersPage(int(page-1), int(count), orderBy, orderAsc)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	tmplVars.Users = *userList

	err = s.templates.ExecuteTemplate(w, "admin_users", tmplVars)
	if err != nil {
		logger.Errorf("could not render dns template: %s", err.Error())
	}
}

func (s *Server) AdminUserAddGetHandler(w http.ResponseWriter, r *http.Request) {
	if u := r.Context().Value(UserKey).(*models.User); !u.IsMemberOfGroup(&models.GroupsUserAdmin) {
		s.returnErrorPage(w, r, http.StatusUnauthorized, "You aren't authorized")
		return
	}

	// Init template variables
	tmplVars := &AdminUsersFormTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmplVars.PageTitle = "Add User"
	tmplVars.Breadcrumbs = &[]templateBreadcrumb{
		{
			Text: "Admin Dashboard",
			HRef: "/app/admin",
		},
		{
			Text: "Users",
			HRef: "/app/admin/users",
		},
		{
			Text: "Add User",
		},
	}

	tmplVars.TitleText = "Add User"
	tmplVars.FormUsername = &templateFormInput{
		ID:          "username",
		Name:        "username",
		Placeholder: "Username",
		Required:    true,
	}
	tmplVars.FormPassword = &templateFormInput{
		ID:          "password",
		Name:        "password",
		Placeholder: "Password",
		Required:    true,
	}
	tmplVars.FormSubmit = &templateFormButton{
		Color: "success",
		Text:  "Add",
	}

	err = s.templates.ExecuteTemplate(w, "admin_user_form", tmplVars)
	if err != nil {
		logger.Errorf("could not render dns template: %s", err.Error())
	}
}

func (s *Server) AdminUserEditGetHandler(w http.ResponseWriter, r *http.Request) {
	if u := r.Context().Value(UserKey).(*models.User); !u.IsMemberOfGroup(&models.GroupsUserAdmin) {
		s.returnErrorPage(w, r, http.StatusUnauthorized, "You aren't authorized")
		return
	}

	// get requested user
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		s.returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}
	user, err := s.db.ReadUser(id)
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if user == nil {
		s.returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("user not found: %s", vars["id"]))
		return
	}

	// Init template variables
	tmplVars := &AdminUsersFormTemplate{}
	err = initTemplate(w, r, tmplVars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmplVars.PageTitle = fmt.Sprintf("Edit User %s", user.Username)
	tmplVars.Breadcrumbs = &[]templateBreadcrumb{
		{
			Text: "Admin Dashboard",
			HRef: "/app/admin",
		},
		{
			Text: "Users",
			HRef: "/app/admin/users",
		},
		{
			Text: fmt.Sprintf("Edit User %s", user.Username),
		},
	}

	tmplVars.TitleText = fmt.Sprintf("Edit User %s", user.Username)
	tmplVars.FormId = &templateFormInput{
		ID:          "id",
		Name:        "id",
		Value:       user.ID.String(),
		Placeholder: "Username",
		Disabled:    true,
	}
	tmplVars.FormUsername = &templateFormInput{
		ID:          "username",
		Name:        "username",
		Value:       user.Username,
		Placeholder: "Username",
		Required:    true,
	}
	tmplVars.FormPassword = &templateFormInput{
		ID:          "password",
		Name:        "password",
		Placeholder: "Password",
	}
	tmplVars.FormSubmit = &templateFormButton{
		Color: "success",
		Text:  "Update",
	}

	err = s.templates.ExecuteTemplate(w, "admin_user_form", tmplVars)
	if err != nil {
		logger.Errorf("could not render dns template: %s", err.Error())
	}
}

