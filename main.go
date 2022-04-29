package main

import routes "github.com/krishnakantha1/to-do-list-backend/Routes"

func main() {
	app := &routes.App{}

	app.InitializeRoutes()
	app.AppStart()
}
