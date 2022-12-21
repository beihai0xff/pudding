package swagger

import (
	"embed"
	"io/fs"
	"mime"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/third_party"
)

// RegisterHandler register Swagger UI handler.
func RegisterHandler(gwmux *runtime.ServeMux, prefix string) {
	swaggerHandler := getOpenAPIHandler(third_party.Embed)
	gwmux.HandlePath("GET", prefix+"/*", func(w http.ResponseWriter,
		r *http.Request, pathParams map[string]string) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
		log.Infof("Serving swagger file: %s", r.URL.Path)
		swaggerHandler.ServeHTTP(w, r)
	})
}

// GetOpenAPIHandler serves an OpenAPI UI.
func getOpenAPIHandler(efs embed.FS) http.Handler {
	mime.AddExtensionType(".svg", "image/svg+xml")

	// Use subdirectory in embedded files
	subFS, err := fs.Sub(efs, "OpenAPI")
	if err != nil {
		panic("couldn't create sub filesystem: " + err.Error())
	}

	return http.FileServer(http.FS(subFS))
}
