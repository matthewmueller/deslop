package naming

import "testing"

func Test_BadName(t *testing.T) { // want `avoid underscores in test names; use "TestBadName" instead of "Test_BadName"`
}

func Test_Another_Bad(t *testing.T) { // want `avoid underscores in test names; use "TestAnother_Bad" instead of "Test_Another_Bad"`
}

func TestGoodName(t *testing.T) {
}
