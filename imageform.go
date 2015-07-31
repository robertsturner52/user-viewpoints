package viewpoints

import (
	"net/http"

	"github.com/nu7hatch/gouuid"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

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
	username := session.Username
	if username == "" {
		http.Redirect(res, req, "/loginform", 302)
		return
	}
	if req.Method == "POST" {
		title := req.FormValue("title")
		caption := req.FormValue("caption")
		imageFile, _, err := req.FormFile("image")
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		defer imageFile.Close()
		imageID, _ := uuid.NewV4()
		err = putFile(ctx, imageID.String(), imageFile)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		image := &Image{
			Username: username,
			Title:    title,
			Caption:  caption,
			ImageID:  imageID.String(),
		}
		key := datastore.NewKey(ctx, "Image", image.ImageID, 0, nil)
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
		Title:    "title",
		ImageID:  "imageID",
		Caption:  "caption",
	}
	err := tpls.ExecuteTemplate(res, "imageform", model)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
}
