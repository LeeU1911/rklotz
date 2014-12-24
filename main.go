package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"./cfg"
	"./controller"
	"./model"
)

func main() {
	defer model.GetDB().Close()

	if err := model.RebuildIndex(); err != nil {
		panic(err)
	}

	if !cfg.GetOptions().Rebuild {
		r := gin.Default()
		r.Use(cfg.Pongo2())

		r.GET("/", controller.FrontController)
		r.GET("/tag/:tag", controller.TagController)
		r.GET("/autocomplete", controller.AutoComplete)

		authorized := r.Group("/@", gin.BasicAuth(gin.Accounts{
			cfg.String("auth.name"): cfg.String("auth.password"),
		}))
		authorized.GET("/new", controller.FormController)
		authorized.POST("/new", controller.FormController)
		authorized.GET("/edit/:uuid", controller.FormController)
		authorized.POST("/edit/:uuid", controller.FormController)
		authorized.GET("/drafts", controller.DraftsController)
		authorized.GET("/published", controller.PublishedController)

		r.NoRoute(controller.PostController)

		r.Static("/static", fmt.Sprintf("%s/static", cfg.GetRootDir()))

		addr := cfg.String("addr")
		cfg.Log(fmt.Sprintf("Running @ %s", addr))
		r.Run(addr)
	}
}
