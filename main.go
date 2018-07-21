package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	b, err := NewBot(os.Getenv("BOT_TOKEN"), os.Getenv("CLIENT_ID"), os.Getenv("ENV"))
	if err != nil {
		panic(err)
	}
	defer b.Close()
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	log.Println("Closing sessions.")
}
