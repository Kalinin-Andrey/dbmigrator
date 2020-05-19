package app


type Logger interface {
	Print(v ...interface{})
	Fatal(v ...interface{})
}



