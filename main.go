package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
)

//go:embed public
var content embed.FS

var (
	Version  = "0.0.1"
	HostPort = flag.String("hostport", "127.0.0.1:8080", "host and port to listen on")
)

func staticFilesEngine() *gin.Engine {
	fsys := fs.FS(content)
	contentStatic, _ := fs.Sub(fsys, "public")

	eng := gin.New()
	eng.Use(gin.Recovery())
	eng.Use(ginzerolog.Logger("gin"))

	eng.StaticFS("/", http.FS(contentStatic))

	return eng
}

func apiEngine() *gin.Engine {
	eng := gin.New()
	eng.Use(gin.Recovery())
	eng.Use(ginzerolog.Logger("gin"))

	apiG := eng.Group("/api")
	{
		apiG.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
	}

	return eng
}

func setupRouter() *gin.Engine {
	api := apiEngine()
	static := staticFilesEngine()

	r := gin.New()
	r.SetTrustedProxies([]string{"0.0.0.0"})
	r.Use(gin.Recovery())
	r.Use(ginzerolog.Logger("gin"))
	r.Any("/*any", func(c *gin.Context) {
		path := c.Param("any")
		if strings.HasPrefix(path, "/api") {
			api.HandleContext(c)
		} else if c.Request.Method == http.MethodGet {
			static.HandleContext(c)
		}
	})

	return r
}

func main() {

	gin.SetMode(gin.ReleaseMode)

	flag.Parse()

	if *HostPort == "" {
		*HostPort = "127.0.0.1:8080"
	}

	r := setupRouter()
	r.Run(fmt.Sprint(*HostPort))
}
