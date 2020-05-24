package app

// Logger interface for application
type Logger interface {
	Print(v ...interface{})
	Fatal(v ...interface{})
}



