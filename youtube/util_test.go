package youtube

import (
	"testing"
)

func TestSafeFileName(t *testing.T) {
	origin := `\/:*?"<>|.`
	safe := safeFileName(origin)
	if safe != "" {
		t.Errorf(`'safeFileName' is not getting rid of the right characters. Expected: "", got: "%s"`, safe)
	}
}
