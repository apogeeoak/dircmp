package compare_test

import (
	"testing"

	"github.com/apogeeoak/dircmp/compare"
)

func TestParseConfigSucceeds(t *testing.T) {
	original := "../test/original"
	compared := "../test/compared"

	config := compare.ParseConfigArgs("", []string{original, compared})

	if config.Original != original {
		t.Fatalf("config.Original improperly set. Wanted '%s'. Got '%s'.", original, config.Original)
	}
	if config.Compared != compared {
		t.Fatalf("config.Compared improperly set. Wanted '%s'. Got '%s'.", compared, config.Compared)
	}
}
