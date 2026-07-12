package example

func init() { // want "avoid init\\(\\) functions; use explicit initialization in constructors or main instead"
	_ = "setup"
}

func Setup() {
	_ = "explicit setup"
}
