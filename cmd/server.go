package cmd

import (
	"goflylivechat/middleware"
	"goflylivechat/router"
	"goflylivechat/tools"
	"goflylivechat/ws"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var (
	port   string
	daemon bool
)

var serverCmd = &cobra.Command{
	Use:     "server",
	Short:   "Start HTTP service",
	Example: "gochat server -p 8080",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	serverCmd.PersistentFlags().StringVarP(&port, "port", "p", "8080", "Port to listen on")
}

func run() {

	baseServer := "0.0.0.0:" + port
	log.Println("Starting server...\nURL: http://" + baseServer)

	// Gin engine setup
	engine := gin.Default()
	engine.LoadHTMLGlob("static/templates/*")
	engine.Static("/static", "./static")
	engine.Use(middleware.SessionHandler())
	engine.Use(middleware.CrossSite)

	// Middlewares
	// engine.Use(middleware.NewMidLogger())

	// Routers
	router.InitViewRouter(engine)
	router.InitApiRouter(engine)

	// Background services
	tools.NewLimitQueue()
	ws.CleanVisitorExpire()
	go ws.WsServerBackend()

	//内置定时任务
	go StartCronJobs()

	// Start server
	engine.Run(baseServer)
}
