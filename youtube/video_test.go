package youtube

import (
	"testing"
)

func TestNewVideo(t *testing.T) {
	v, err := NewVideo("Nq5LMGtBmis")
	if err != nil {
		t.Error(err)
	}
	if v.ID != "Nq5LMGtBmis" {
		t.Error("wrong id")
	}
	if v.Author != "Vulf" {
		t.Error("wrong author")
	}
}
