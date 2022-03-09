package main

import (
	"github.com/sevings/zhk_bot/internal"
	"log"
	"os"
	"os/signal"
)

func main() {
	bot := internal.NewBot()

	go bot.Run()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown bot")
	bot.Stop()
}
