package swagger

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path"

	"github.com/go-keg/swagger-api/dist"
	"gopkg.in/yaml.v3"
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
	urls          []OpenapiURL
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

func SetURLs(urls []OpenapiURL) Option {
	return func(cfg *config) {
		for i := range urls {
			urls[i].URL = path.Join(cfg.OpenapiPath(), urls[i].URL)
		}
		cfg.urls = urls
	}
}

type Openapi struct {
	Info struct {
		Title string
	}
}

// Handler swagger ui
func Handler(apis fs.FS, opts ...Option) http.Handler {
	cfg := config{
		prefix:        "swagger",
		swaggerUIPath: "ui",
		openapiPath:   "apis",
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	if len(cfg.urls) == 0 {
		err := fs.WalkDir(apis, ".", func(filePath string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				content, err := apis.Open(filePath)
				if err != nil {
					return err
				}
				var openapi Openapi
				decoder := yaml.NewDecoder(content)
				if err := decoder.Decode(&openapi); err != nil {
					panic(fmt.Errorf("解析配置文件失败: %v", err))
				}
				cfg.urls = append(cfg.urls, OpenapiURL{
					URL:  path.Join(cfg.OpenapiPath(), filePath),
					Name: openapi.Info.Title,
				})
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	}

	apiServer := http.FileServer(http.FS(apis))
	staticServer := http.FileServer(http.FS(dist.SwagFS))
	mux := http.NewServeMux()
	mux.Handle(cfg.OpenapiPath()+"/", http.StripPrefix(cfg.OpenapiPath(), apiServer))
	mux.Handle(cfg.SwaggerUIPath()+"/public/", http.StripPrefix(cfg.SwaggerUIPath()+"/public", staticServer))
	mux.Handle(cfg.SwaggerUIPath(), renderIndex(map[string]any{
		"URLs":   cfg.urls,
		"Prefix": cfg.SwaggerUIPath() + "/public",
	}))
	return mux
}

func renderIndex(data any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("swagger-ui").Parse(IndexTemp)
		if err != nil {
			http.Error(w, "Could not load template", http.StatusInternalServerError)
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Could not render template", http.StatusInternalServerError)
		}
	}
}

type OpenapiURL struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
