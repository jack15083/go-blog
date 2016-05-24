package main

import (
	log "github.com/Sirupsen/logrus"

	"github.com/julienschmidt/httprouter"
	"net/http"

	"github.com/beforydeath/go-blog/controllers"
	"github.com/beforydeath/go-blog/core"
	"github.com/beforydeath/go-blog/models"
)

func main() {
	core.Config.Init()
	core.Themes.Init()

	core.Config.Logger.Init()
	if core.Config.Logger.OutFile {
		logFile := core.Config.Logger.File()
		defer logFile.Close()
	}

	models.ConnectDB()
	defer models.CloseDB()

	c := new(controllers.Controller)

	c.User.Name = []byte("admin")
	c.User.Password = []byte("admin")

	router := httprouter.New()

	router.ServeFiles("/vendor/*filepath", http.FileSystem(http.Dir(core.Config.BasePath + "/themes/assets/vendors/")))
	router.ServeFiles("/assets/*filepath", http.FileSystem(http.Dir(core.Config.ThemePath + "/assets/")))

	router.GET("/", c.Page.GetList)
	router.GET("/page/:alias", c.Page.GetByAlias)

	router.GET("/login", c.User.BasicAuth(c.Page.GetList))
	router.GET("/admin/page/create", c.User.BasicAuth(c.Page.Create))
	router.GET("/admin/page/update/:id", c.User.BasicAuth(c.Page.Update))
	router.POST("/admin/page/create", c.User.BasicAuth(c.Page.Create))
	router.POST("/admin/page/update/:id", c.User.BasicAuth(c.Page.Update))
	router.GET("/admin/page/delete/:id", c.User.BasicAuth(c.Page.Delete))

	log.Info("Server started ...")
	log.Fatal(http.ListenAndServe(":8085", router))
}
