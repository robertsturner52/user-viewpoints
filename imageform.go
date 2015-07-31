package viewpoints

import (
	"net/http"
	"strings"

	"github.com/nu7hatch/gouuid"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type Image struct {
	Username string
	Title    string
	ImageID  string
	Caption  string
}

type imageModel struct {
	Session  Session
	ImageID  string
	Username string
	Title    string
	Caption  string
}

func handleImageForm(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	session := getSession(ctx, req)
	imageID := strings.SplitN(req.URL.Path, "/", 3)[2]
	username := session.Username
	if username == "" {
		http.Redirect(res, req, "/loginform", 302)
		return
	}
	image := &Image{
		Username: username,
	}
	var key *datastore.Key
	if imageID != "" {
		key = datastore.NewKey(ctx, "Image", imageID, 0, nil)
		err := datastore.Get(ctx, key, image)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
	} else {
		id, _ := uuid.NewV4()
		key = datastore.NewKey(ctx, "Image", id.String(), 0, nil)
		image.ImageID = id.String()
	}

	if req.Method == "POST" {
		image.Title = req.FormValue("title")
		image.Caption = req.FormValue("caption")
		imageFile, _, err := req.FormFile("image")
		if err == nil {
			defer imageFile.Close()
			err = putFile(ctx, image.ImageID, imageFile)
			if err != nil {
				http.Error(res, err.Error(), 500)
				return
			}
		} else if imageID == "" {
			http.Error(res, err.Error(), 500)
			return
		}
		_, err = datastore.Put(ctx, key, image)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		http.Redirect(res, req, "/user/"+username, 302)
		return
	}
	model := &imageModel{
		Username: session.Username,
		Title:    image.Title,
		ImageID:  image.ImageID,
		Caption:  image.Caption,
	}
	err := tpls.ExecuteTemplate(res, "imageform", model)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
}
