package viewpoints

import (
	"net/http"
	"text/template"

	"google.golang.org/appengine"
)

func init() {
	http.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir("images/"))))
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("css/"))))
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/user/", handlePersonalPage)
	http.HandleFunc("/userslist", handleUsersList)
	http.HandleFunc("/imageform/", handleImageForm)
	http.HandleFunc("/listform", handleListForm)
	http.HandleFunc("/register", handleRegistration)
	http.HandleFunc("/loginform", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/personal-images/", handlePersonalImages)
}

var tpls = template.Must(template.ParseGlob("templates/*.gohtml"))

type indexModel struct {
	Session Session
}

func handleIndex(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}
	ctx := appengine.NewContext(req)
	session := getSession(ctx, req)
	model := &indexModel{
		Session: session,
	}
	err := tpls.ExecuteTemplate(res, "homepage", model)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
}
