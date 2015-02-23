package logging

type LogInterface interface {
	Info(tag, data string)  // used to start a container starting (it does all the operations, like pull and start)
	Error(tag, data string) // stops the container and removes the stopped container
}
