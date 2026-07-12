package example

import "errors"

func bad() error {
	err := doThing()
	if err == nil { // want `avoid "if err == nil"`
		return nil
	}
	return err
}

// Allowed: "nil == err" is an intentional escape hatch.
func allowed() error {
	err := doThing()
	if nil == err {
		return nil
	}
	return err
}

// Allowed: normal err != nil with no else.
func good() error {
	err := doThing()
	if err != nil {
		return err
	}
	return nil
}

func doThing() error { return errors.New("oops") }
