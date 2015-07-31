package viewpoints

import (
	"net/http"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type personalModel struct {
	Session   Session
	Username  string
	ListItems []string
	Images    []Image
}

func handlePersonalImages(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	name := strings.SplitN(req.URL.Path, "/", 3)[2]
	link, err := getFileLink(ctx, name)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	http.Redirect(res, req, link, 301)
}

func handlePersonalPage(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	session := getSession(ctx, req)
	username := strings.SplitN(req.URL.Path, "/", 3)[2]
	if username == "" {
		http.Redirect(res, req, "/userslist", 302)
		return
	}
	q := datastore.NewQuery("Image").Filter("Username =", username)
	var images []Image
	q.GetAll(ctx, &images)
	key := datastore.NewKey(ctx, "List", username, 0, nil)
	var list List
	datastore.Get(ctx, key, &list)
	listItems := make([]string, 0)
	if strings.TrimSpace(list.List) != "" {
		listItems = strings.Split(list.List, "\n")
	}
	model := &personalModel{
		Session:   session,
		Username:  username,
		ListItems: listItems,
		Images:    images,
	}
	err := tpls.ExecuteTemplate(res, "personalpage", model)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
}
