package batch

import "testing"

func TestRemoveAfterLastNL(t *testing.T) {
	var dataString = "test test\n test"
	var expectedString = "test test"
	var returnedString string
	var expectedToCut = int64(6)
	var returnedToCut int64
	returnedString, returnedToCut = removeAfterLastNL(dataString)
	if returnedString != expectedString {
		t.Error("Expected ", expectedString, ", got ", returnedString)
	}

	if returnedToCut != expectedToCut {
		t.Error("Expected ", expectedToCut, ", got ", returnedToCut)
	}
}
