package main

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/fujiwara/logutils"
)

func main() {
	filter := &logutils.LevelFilter{
		Levels: []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
		ModifierFuncs: []logutils.ModifierFunc{
			nil, // default
			logutils.Color(color.FgYellow),
			logutils.Color(color.FgRed, color.BgBlack),
		},
		MinLevel: logutils.LogLevel("WARN"),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)

	log.Print("[DEBUG] Debugging")         // this will not print
	log.Print("[WARN] Warning")            // this will print as yellow font
	log.Print("[ERROR] Erring")            // this will print as red font and black background
	log.Print("Message I haven't updated") // and so will this
}
