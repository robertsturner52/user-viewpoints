package viewpoints

import (
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"
)

const bucketName = "users-viewpoints-list.appspot.com"

func getCloudContext(appengineContext context.Context) context.Context {
	hc := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(appengineContext, storage.ScopeFullControl),
			Base:   &urlfetch.Transport{Context: appengineContext},
		},
	}

	return cloud.NewContext(appengine.AppID(appengineContext), hc)
}

func putFile(ctx context.Context, name string, rdr io.Reader) error {
	cctx := getCloudContext(ctx)
	writer := storage.NewWriter(cctx, bucketName, name)
	writer.ACL = []storage.ACLRule{
		{storage.AllUsers, storage.RoleReader},
	}
	io.Copy(writer, rdr)
	return writer.Close()
}

func getFile(ctx context.Context, name string) (io.ReadCloser, error) {
	cctx := getCloudContext(ctx)
	return storage.NewReader(cctx, bucketName, name)
}

func getFileLink(ctx context.Context, name string) (string, error) {
	cctx := getCloudContext(ctx)
	obj, err := storage.StatObject(cctx, bucketName, name)
	if err != nil {
		return "", err
	}
	return obj.MediaLink, nil
}

func handleImages(res http.ResponseWriter, req *http.Request) {
	imageID := strings.SplitN(req.URL.Path, "/", 3)[2]
	ctx := appengine.NewContext(req)
	mediaLink, err := getFileLink(ctx, imageID)
	if err != nil {
		http.NotFound(res, req)
		return
	}
	http.Redirect(res, req, mediaLink, 302)
}
