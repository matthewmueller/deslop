package shadow

// OK: log is a local variable, not the stdlib log package.
func Good() {
	log := &logger{}
	log.Println("hello")
}

type logger struct{}

func (l *logger) Println(args ...any) {}
