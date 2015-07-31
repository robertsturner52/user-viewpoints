package viewpoints

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type Profile struct {
	Username     string
	PasswordHash string
}
type registerModel struct {
	Session Session
	Error   string
}
type loginModel struct {
	Error   string
	Session Session
}

func (p *Profile) setPassword(ctx context.Context, password string) {
	bs, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	p.PasswordHash = fmt.Sprintf("%x", bs)
}

func (p *Profile) checkPassword(ctx context.Context, password string) bool {
	bs, _ := hex.DecodeString(p.PasswordHash)
	err := bcrypt.CompareHashAndPassword(bs, []byte(password))
	return err == nil
}

func createProfile(ctx context.Context, profile *Profile) error {
	key := datastore.NewKey(ctx, "Profile", profile.Username, 0, nil)
	_, err := datastore.Put(ctx, key, profile)
	return err
}

func handleRegistration(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	session := getSession(ctx, req)
	model := &registerModel{}
	if req.Method == "POST" {
		username := req.FormValue("username")
		password := req.FormValue("password")
		confirm := req.FormValue("confirm")
		if password == "" {
			model.Error = "Password is required."
		} else if password != confirm {
			model.Error = "Passwords must match."
		} else {
			var profile Profile
			key := datastore.NewKey(ctx, "Profile", username, 0, nil)
			datastore.Get(ctx, key, &profile)
			if username == profile.Username {
				model.Error = "Username already taken."
			} else {
				profile := &Profile{
					Username: username,
				}
				profile.setPassword(ctx, password)
				createProfile(ctx, profile)
				session.Username = username
				putSession(ctx, res, session)
				http.Redirect(res, req, "/listform", 302)
				return
			}
		}
	}
	err := tpls.ExecuteTemplate(res, "register", model)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
}

func handleLogin(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	session := getSession(ctx, req)
	model := &loginModel{}
	if req.Method == "POST" {
		username := req.FormValue("username")
		password := req.FormValue("password")
		var profile Profile
		key := datastore.NewKey(ctx, "Profile", username, 0, nil)
		datastore.Get(ctx, key, &profile)
		if profile.checkPassword(ctx, password) {
			session.Username = username
			putSession(ctx, res, session)
			http.Redirect(res, req, "user/"+username, 302)
			return
		} else {
			model.Error = "Invalid username or password."
		}
	}
	err := tpls.ExecuteTemplate(res, "loginform", model)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
}

func handleLogout(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	session := getSession(ctx, req)
	session.Username = ""
	putSession(ctx, res, session)
	http.Redirect(res, req, "/", 302)
}
