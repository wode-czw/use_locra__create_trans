package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wode-czw/trans_czw_wode/server/config"
	"github.com/wode-czw/trans_czw_wode/server/controllers"

	"github.com/wode-czw/trans_czw_wode/server/ws"
)

//go:embed frontend/dist/*
var FS embed.FS //这行是说下面你执行的和FS这个变量有关系的这个文件是要被打包的

func Run() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	//hub和gin可以说没有关系
	hub := ws.NewHub()
	go hub.Run()
	router.GET("/ws", func(c *gin.Context) {
		ws.HttpController(c, hub)
	})

	//####
	router.POST("api/v1/texts", controllers.TextsController)

	router.GET("api/v1/addresses", controllers.AddressesController)
	router.POST("/api/v1/files", controllers.FilesController)
	router.GET("/uploads/:path", controllers.UploadsController)

	router.GET("api/v1/qrcodes", controllers.QrcodesController) //可是电脑并没有主动跟你要二维码呀

	staticFiles, _ := fs.Sub(FS, "frontend/dist")
	router.StaticFS("/static", http.FS(staticFiles))

	//##########你的确是static开头，但是你其他语句找不到route
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/static/") {
			reader, err := staticFiles.Open("index.html")
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()
			stat, err := reader.Stat()
			if err != nil {
				log.Fatal(err)
			}
			c.DataFromReader(http.StatusOK, stat.Size(), "text/html", reader, nil)
		} else {
			c.Status(http.StatusNotFound)
		}
	})

	router.Run(config.Get_port())
}
