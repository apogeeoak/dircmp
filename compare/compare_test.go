package compare_test

import (
	"testing"

	"github.com/apogeeoak/dircmp/compare"
)

func TestCompareBasic(t *testing.T) {
	config := setupConfigBasic()
	want := "Searched 8 file(s), 4 file(s) different, 2 director(ies) different, 6 total entr(ies) different"

	stat, err := compare.Compare(config)

	if err != nil {
		t.Fatal("Error:", err)
	}
	if stat.String() != want {
		t.Fatalf("Wanted '%s'. Got '%s'.", want, stat)
	}
}

func setupConfigBasic() *compare.Config {
	original := "../test/original"
	compared := "../test/compared"
	return compare.ParseConfigArgs("", []string{original, compared})
}
