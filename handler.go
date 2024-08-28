package swagger

import (
	"html/template"
	"net/http"
	"path"

	"github.com/eiixy/swagger-api/dist"
	"github.com/gin-gonic/gin"
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
	prefix        string
	swaggerUIPath string
	openapiPath   string
}

func (r config) SwaggerUIPath() string {
	return path.Join("/", r.prefix, r.swaggerUIPath)
}

func (r config) OpenapiPath() string {
	return path.Join("/", r.prefix, r.openapiPath)
}

type Option func(cfg *config)

func SetPrefix(prefix string) Option {
	return func(cfg *config) {
		cfg.prefix = prefix
	}
}

func SetSwaggerUIPath(path string) Option {
	return func(cfg *config) {
		cfg.swaggerUIPath = path
	}
}

func SetOpenapiPath(path string) Option {
	return func(cfg *config) {
		cfg.openapiPath = path
	}
}

// Handler swagger ui
func Handler(apis http.FileSystem, urls []OpenapiURL, opts ...Option) http.Handler {
	cfg := config{
		prefix:        "swagger",
		swaggerUIPath: "ui",
		openapiPath:   "apis",
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	for i := range urls {
		urls[i].URL = cfg.OpenapiPath() + urls[i].URL
	}

	router := gin.New()
	router.SetHTMLTemplate(template.Must(template.New("swagger-ui").Parse(IndexTemp)))
	gin.SetMode(gin.ReleaseMode)
	router.GET(cfg.SwaggerUIPath(), func(c *gin.Context) {
		c.HTML(200, "swagger-ui", map[string]any{
			"URLs":   urls,
			"Prefix": cfg.SwaggerUIPath() + "/public",
		})
	})
	router.StaticFS(cfg.SwaggerUIPath()+"/public", http.FS(dist.SwagFS))
	router.StaticFS(cfg.OpenapiPath(), apis)
	return router
}

type OpenapiURL struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
