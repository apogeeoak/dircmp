package compare_test

import (
	"testing"

	"github.com/apogeeoak/dircmp/compare"
)

func TestParseConfigBasic(t *testing.T) {
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

func TestParseConfigEntire(t *testing.T) {
	entire := "--entire"
	dir := "../test"
	want := int64(0)

	config := compare.ParseConfigArgs("", []string{entire, dir, dir})
	offset := config.Offset(1)

	if offset != want {
		t.Fatalf("config.Offset improperly set. Wanted '%v'. Got '%v'.", want, offset)
	}
}
