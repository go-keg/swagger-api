//go:generate statik -src=../../zapis/api -include=*.openapi.yaml -ns apis -p apis
package swagger

import (
	_ "github.com/eiixy/swagger-api/swagger"
	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"
	"html/template"
	"net/http"
)

const IndexTemp = `<!-- HTML for static distribution bundle build -->
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="{{.Prefix}}/swagger-ui.css"/>
    <link rel="stylesheet" type="text/css" href="{{.Prefix}}/index.css"/>
    <link rel="icon" type="image/png" href="{{.Prefix}}/favicon-32x32.png" sizes="32x32"/>
    <link rel="icon" type="image/png" href=".{{.Prefix}}/favicon-16x16.png" sizes="16x16"/>
</head>

<body>
<div id="swagger-ui"></div>
<script src="{{.Prefix}}/swagger-ui-bundle.js" charset="UTF-8"></script>
<script src="{{.Prefix}}/swagger-ui-standalone-preset.js" charset="UTF-8"></script>
<script>
    window.onload = function () {
        window.ui = SwaggerUIBundle({
            urls: {{.URLs}},
            dom_id: '#swagger-ui',
            deepLinking: true,
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIStandalonePreset
            ],
            plugins: [
                SwaggerUIBundle.plugins.DownloadUrl
            ],
            layout: "StandaloneLayout"
        });
    };
</script>
</body>
</html>`

type config struct {
	swaggerUIPrefix string
	openapiPrefix   string
}

type Option func(cfg *config)

func SetSwaggerUIPrefix(prefix string) Option {
	return func(cfg *config) {
		cfg.swaggerUIPrefix = prefix
	}
}

func SetOpenapiPrefix(prefix string) Option {
	return func(cfg *config) {
		cfg.openapiPrefix = prefix
	}
}

// Handler swagger ui
func Handler(apis http.FileSystem, urls []OpenapiURL, opts ...Option) http.Handler {
	cfg := config{
		swaggerUIPrefix: "/swagger-ui",
		openapiPrefix:   "/openapi/apis",
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	for i := range urls {
		urls[i].URL = cfg.openapiPrefix + urls[i].URL
	}
	swaggerFS, err := fs.NewWithNamespace("swagger-ui")
	if err != nil {
		panic(err)
	}

	router := gin.New()
	router.SetHTMLTemplate(template.Must(template.New("swagger-ui").Parse(IndexTemp)))
	gin.SetMode(gin.ReleaseMode)
	router.GET(cfg.swaggerUIPrefix, func(c *gin.Context) {
		c.HTML(200, "swagger-ui", map[string]any{
			"URLs":   urls,
			"Prefix": cfg.swaggerUIPrefix + "/public",
		})
	})
	router.StaticFS(cfg.swaggerUIPrefix+"/public", swaggerFS)
	router.StaticFS(cfg.openapiPrefix, apis)
	return router
}

type OpenapiURL struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type fileSystem struct {
	fs        http.FileSystem
	indexFile http.File
}

func (r fileSystem) Open(name string) (http.File, error) {
	f, err := r.fs.Open(name)
	if err != nil {
		return nil, err
	}
	if name == "/index.html" {
		return r.indexFile, nil
	}
	return f, err
}
