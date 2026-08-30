package anytype

import (
	"encoding/json"
	"testing"
)

func TestNewPropertyLinkValue(t *testing.T) {
	cases := []struct {
		format PropertyFormat
		value  any
		want   string
	}{
		{PropertyFormatText, "hi", `{"key":"k","text":"hi"}`},
		{PropertyFormatNumber, 3, `{"key":"k","number":3}`},
		{PropertyFormatNumber, 2.5, `{"key":"k","number":2.5}`},
		{PropertyFormatCheckbox, true, `{"key":"k","checkbox":true}`},
		{PropertyFormatSelect, "tag", `{"key":"k","select":"tag"}`},
		{PropertyFormatSelect, nil, `{"key":"k","select":null}`},
		{PropertyFormatDate, nil, `{"key":"k","date":null}`},
		{PropertyFormatMultiSelect, []string{"a", "b"}, `{"key":"k","multi_select":["a","b"]}`},
		{PropertyFormatFiles, []string{"f"}, `{"key":"k","files":["f"]}`},
		{PropertyFormatObjects, []string{"o"}, `{"key":"k","objects":["o"]}`},
		{PropertyFormatURL, "https://x", `{"key":"k","url":"https://x"}`},
		{PropertyFormatEmail, "a@b", `{"key":"k","email":"a@b"}`},
		{PropertyFormatPhone, "+1", `{"key":"k","phone":"+1"}`},
	}
	for _, c := range cases {
		v, err := NewPropertyLinkValue("k", c.format, c.value)
		if err != nil {
			t.Fatalf("%s: %v", c.format, err)
		}
		got, _ := json.Marshal(v)
		if string(got) != c.want {
			t.Errorf("%s: got %s, want %s", c.format, got, c.want)
		}
	}
}

func TestNewPropertyLinkValueRejects(t *testing.T) {
	bad := []struct {
		format PropertyFormat
		value  any
	}{
		{PropertyFormatText, 1},
		{PropertyFormatNumber, "1"},
		{PropertyFormatCheckbox, "true"},
		{PropertyFormatMultiSelect, "a"},
		{PropertyFormatSelect, 1},
		{PropertyFormat("color"), "red"},
	}
	for _, c := range bad {
		if _, err := NewPropertyLinkValue("k", c.format, c.value); err == nil {
			t.Errorf("%s with %T: want error", c.format, c.value)
		}
	}
}

func TestIconHelpers(t *testing.T) {
	got, _ := json.Marshal(NamedIcon("star", ColorRed))
	if string(got) != `{"format":"icon","name":"star","color":"red"}` {
		t.Errorf("got %s", got)
	}
	got, _ = json.Marshal(EmojiIcon("x"))
	if string(got) != `{"format":"emoji","emoji":"x"}` {
		t.Errorf("got %s", got)
	}
}
