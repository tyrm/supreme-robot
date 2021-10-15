package webapp

import (
	"github.com/tyrm/supreme-robot/models"
	"github.com/tyrm/supreme-robot/util"
	"net/http"
	"regexp"
)

func makeNavbar(r *http.Request) (navbar *[]templateNavbarNode) {

	// create navbar
	newNavbar := []templateNavbarNode{
		{
			Text:     "Home",
			MatchStr: "^/app/$",
			FAIcon:   "home",
			URL:      "/app/",
		},
		{
			Text:     "DNS",
			MatchStr: "^/app/dns",
			FAIcon:   "map",
			URL:      "/app/dns",
		},
	}

	// Show Admin Menu
	if r.Context().Value(UserKey) != nil {
		user := r.Context().Value(UserKey).(*models.User)

		if util.ContainsOneOfStrings(user.Groups, models.GroupsAllAdmins) {

			adminMenu := templateNavbarNode{
				Text:     "Admin",
				MatchStr: "^/app/admin",
				FAIcon:   "hammer",
				URL:      "#",
			}

			if util.ContainsOneOfStrings(user.Groups, models.GroupsUserAdmin) {
				adminMenu.Children = append(adminMenu.Children, &templateNavbarNode{
					Text:     "Users",
					MatchStr: "^/app/admin/users",
					FAIcon:   "users",
					URL:      "/app/admin/users",
				})
			}

			newNavbar = append(newNavbar, adminMenu)
		}
	}

	for i := 0; i < len(newNavbar); i++ {
		if newNavbar[i].MatchStr != "" {
			match, err := regexp.MatchString(newNavbar[i].MatchStr, r.URL.Path)
			if err != nil {
				logger.Errorf("makeNavbar:Error matching regex: %v", err)
			}
			if match {
				newNavbar[i].Active = true
			}

		}

		if newNavbar[i].Children != nil {
			for j := 0; j < len(newNavbar[i].Children); j++ {

				if newNavbar[i].Children[j].MatchStr != "" {
					subMatch, err := regexp.MatchString(newNavbar[i].Children[j].MatchStr, r.URL.Path)
					if err != nil {
						logger.Errorf("makeNavbar:Error matching regex: %v", err)
					}

					if subMatch {
						newNavbar[i].Active = true
						newNavbar[i].Children[j].Active = true
					}

				}

			}
		}
	}

	return &newNavbar
}
