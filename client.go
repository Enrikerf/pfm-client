package main

import "github.com/Enrikerf/pfm/commandExecutor/app/Config"

func main() {
	var app = Config.App{}
	app.Run()

	//var appEngine = Config.NewEngineApp()
	//appEngine.Run()
	//select {}
}
