package viewpoints

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type userslistModel struct {
	Session  Session
	Profiles []Profile
}

func handleUsersList(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	q := datastore.NewQuery("Profile").Order("Username")
	var profiles []Profile
	_, err := q.GetAll(ctx, &profiles)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	model := &userslistModel{
		Profiles: profiles,
	}
	err = tpls.ExecuteTemplate(res, "userslist", model)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
}
