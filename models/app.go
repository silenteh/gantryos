package models

// this file is the main one to deploy apps
// this has the concept of deployments and groups and amount of instances

// import (
// 	""
// )

type App struct {
	Id        string // unique auto generated id
	Name      string // the application name - just a reference for users to understand which app it is
	Version   string // the application version
	Instances int    // amount of instances to run
	Task      Task   // the task information to run the app/cmd/vm/tasks
}
