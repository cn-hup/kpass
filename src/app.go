package src

import (
	"bytes"
	"regexp"
	"strings"
	"time"

	"github.com/seccom/kpass/src/auth"
	"github.com/seccom/kpass/src/logger"
	"github.com/seccom/kpass/src/service"
	"github.com/teambition/gear"
	"github.com/teambition/gear/middleware/cors"
	"github.com/teambition/gear/middleware/favicon"
	"github.com/teambition/gear/middleware/secure"
	"github.com/teambition/gear/middleware/static"
	"github.com/tidwall/buntdb"
)

// Version is app version
const Version = "v1.0.0-alpha.5"

// New returns a app instance
func New(dbPath, bindHost string) (*gear.App, *buntdb.DB) {
	app := gear.New()
	if app.Env() == "production" {
		logger.Init()
	}

	db, err := service.NewDB(dbPath)
	if err != nil {
		panic(err)
	}
	auth.Init(db.Salt, 20*time.Minute)

	indexBody := "<h1>Kpass</h1>"
	faviconBin := []byte{}

	if app.Env() != "test" {
		indexBody = string(MustAsset("index.html"))
		faviconBin = MustAsset("favicon.ico")
	}

	staticOpts := static.Options{
		Root:        "",
		Prefix:      "/static/",
		StripPrefix: false,
		Files:       make(map[string][]byte),
		Includes:    []string{"/logo.png", "/humans.txt", "/robots.txt", "/kpass.png"},
	}
	for _, name := range AssetNames() {
		name = "/" + name
		staticOpts.Files[name] = MustAsset(name[1:])
		if bindHost != "" && strings.HasSuffix(name, ".js") {
			staticOpts.Files[name] = bytes.Replace(staticOpts.Files[name],
				[]byte("http://127.0.0.1:8088"), []byte(bindHost), -1)
		}
	}
	if app.Env() == "development" {
		staticOpts.Root = "./web"
	}

	app.Use(favicon.NewWithIco(faviconBin))
	app.Use(static.New(staticOpts))
	app.UseHandler(logger.Default())

	router := newRouter(db)
	router.Use(cors.New())
	router.Use(secure.Default)
	routerPrefix := regexp.MustCompile(`^/(api|download|upload)/`)
	router.Otherwise(func(ctx *gear.Context) (err error) {
		if routerPrefix.MatchString(ctx.Path) {
			return ctx.ErrorStatus(404)
		}
		return ctx.HTML(200, indexBody)
	})
	app.UseHandler(router)

	return app, db.DB
}
