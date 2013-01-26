package mango

import (
	"testing"
)

func TestMimeType(t *testing.T) {
	test := func(value, expected string) {
		found := MimeType(value, "")
		if found != expected {
			t.Error("Expected", value, "to have mime type:", expected, "got:", found)
		}
	}

	test(".css", "text/css")
	test(".js", "application/javascript")
}

func TestMimeTypeFallback(t *testing.T) {
	value := "bogus"
	expected := "fallback/value"
	found := MimeType(value, expected)
	if found != expected {
		t.Error("Expected", value, "to have fallback mime type:", expected, "got:", found)
	}
}

func TestAddingMimeTypes(t *testing.T) {
	value := ".new"
	expected := "new/type"
	fallback := "fallback/type"

	found := MimeType(value, fallback)
	if found != fallback {
		t.Error("Expected", value, "to have fallback mime type:", fallback, "got:", found)
	}

	MimeTypes[".new"] = "new/type"

	found = MimeType(value, fallback)
	if found != expected {
		t.Error("Expected", value, "to have new mime type:", expected, "got:", found)
	}
}

func TestRemovingMimeTypes(t *testing.T) {
	value := ".jpg"
	expected := "image/jpeg"
	fallback := "fallback/type"

	found := MimeType(value, fallback)
	if found != expected {
		t.Error("Expected", value, "to have mime type:", expected, "got:", found)
	}

	delete(MimeTypes, ".jpg")

	found = MimeType(value, fallback)
	if found != fallback {
		t.Error("Expected", value, "to have fallback mime type:", fallback, "got:", found)
	}
}
