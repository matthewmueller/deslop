package example

import "os"

func Bad() {
	os.Getenv("FOO")    // want "do not use os.Getenv directly; add FOO to the Env struct"
	os.LookupEnv("BAR") // want "do not use os.LookupEnv directly; add BAR to the Env struct"
}
