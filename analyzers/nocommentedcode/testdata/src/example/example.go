package example

// Doc comment for Foo is fine.
func Foo() {
	// TODO: implement this later
	// NOTE: this is important
	// nolint:errcheck

	// x := doSomething() // want "commented-out code should be deleted, not commented; use version control to recover old code"
	// fmt.Println("hello") // want "commented-out code should be deleted, not commented; use version control to recover old code"
	// if err != nil { // want "commented-out code should be deleted, not commented; use version control to recover old code"
	// return nil, err // want "commented-out code should be deleted, not commented; use version control to recover old code"
	// server.Close() // want "commented-out code should be deleted, not commented; use version control to recover old code"

	// This is a normal comment explaining something.
	// Another normal explanation line.
	_ = "real code"
}
