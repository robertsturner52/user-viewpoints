package viewpoints

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type List struct {
	Username string
	List     string
}
type listModel struct {
	Session Session
	ListID  string
	List    string
}

func handleListForm(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	session := getSession(ctx, req)
	username := session.Username
	if username == "" {
		http.Redirect(res, req, "/loginform", 302)
		return
	}
	list := &List{
		Username: username,
	}
	key := datastore.NewKey(ctx, "List", username, 0, nil)
	datastore.Get(ctx, key, list)
	if req.Method == "POST" {
		list.List = req.FormValue("list")
		_, err := datastore.Put(ctx, key, list)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		http.Redirect(res, req, "/user/"+username, 302)
		return
	}
	model := &listModel{
		Session: session,
		List:    list.List,
	}
	err := tpls.ExecuteTemplate(res, "listform", model)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
}
