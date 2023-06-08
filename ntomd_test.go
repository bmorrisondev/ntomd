package ntomd

import (
	"testing"

	"github.com/dstotijn/go-notion"
)

func Test_ParseRichText_Standard(t *testing.T) {
	expected := "Hello world"
	actual := ParseRichText([]notion.RichText{
		{
			Text: &notion.Text{
				Content: "Hello ",
			},
		},
		{
			Text: &notion.Text{
				Content: "world",
			},
		},
	})

	if actual == nil {
		t.Fatalf("expected: %v, actual is nil", expected)
	}
	if expected != *actual {
		t.Fatalf("expected: %v, actual: %v", expected, *actual)
	}
}

func Test_ParseRichText_Bold(t *testing.T) {
	expected := "**Hello** world"
	actual := ParseRichText([]notion.RichText{
		{
			Text: &notion.Text{
				Content: "Hello",
			},
			Annotations: &notion.Annotations{
				Bold: true,
			},
		},
		{
			Text: &notion.Text{
				Content: " world",
			},
		},
	})

	if actual == nil {
		t.Fatalf("expected: %v, actual is nil", expected)
	}
	if expected != *actual {
		t.Fatalf("expected: %v, actual: %v", expected, *actual)
	}
}

func Test_ParseRichText_BoldItalic(t *testing.T) {
	expected := "**_Hello world_**"
	actual := ParseRichText([]notion.RichText{
		{
			Text: &notion.Text{
				Content: "Hello world",
			},
			Annotations: &notion.Annotations{
				Bold:   true,
				Italic: true,
			},
		},
	})

	if actual == nil {
		t.Fatalf("expected: %v, actual is nil", expected)
	}
	if expected != *actual {
		t.Fatalf("expected: %v, actual: %v", expected, *actual)
	}
}
