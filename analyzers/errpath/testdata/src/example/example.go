package example

import "errors"

func bad() error {
	err := doThing()
	if err == nil { // want `avoid "if err == nil"`
		// happy path nested...
		return nil
	}
	return err
}

func badElse() error {
	err := doThing()
	if err != nil {
		return err
	} else { // want "unnecessary else after return/continue/break"
		// happy path in else...
		return nil
	}
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
	// happy path on the left...
	return nil
}

func doThing() error { return errors.New("oops") }
