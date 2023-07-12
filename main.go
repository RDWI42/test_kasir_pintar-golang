package main

import (
	"belajar-golang/controller"
	"os"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Inisialisasi log
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("api.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	e := echo.New()
	e.POST("/regis", controller.Regis)
	e.POST("/inquiry", controller.Inquiry)
	e.POST("/payment", controller.Payment)
	e.Logger.Fatal(e.Start(":8000"))
}
