package compare_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/apogeeoak/dircmp/compare"
	"github.com/apogeeoak/dircmp/lib/test"
)

func TestCompareBasic(t *testing.T) {
	config := setupConfigBasic()
	want := "Searched 8 file(s), 4 file(s) different, 2 director(ies) different, 6 total entr(ies) different. 0 error(s)."

	stats, err := compare.Compare(config)
	fmt.Println(stats)

	if err != nil {
		t.Fatal("Error:", err)
	}
	if stats.String() != want {
		t.Fatalf("Wanted '%s'. Got '%s'.", want, stats)
	}
}

func TestCompareLarge(t *testing.T) {
	config := setupConfigLarge()
	want := compareSameRegex

	stats, err := compare.Compare(config)
	fmt.Println(stats)

	if err != nil {
		t.Fatal("Error:", err)
	}
	if !want.MatchString(stats.String()) {
		t.Fatalf("Wanted '%s'. Got '%s'.", want, stats)
	}
}

func TestCompareEntire(t *testing.T) {
	// Require the long flag to be set in order to run this long running test.
	test.RequireLong(t)

	config := setupConfigEntire()
	want := compareSameRegex

	stats, err := compare.Compare(config)
	fmt.Println(stats)

	if err != nil {
		t.Fatal("Error:", err)
	}
	if !want.MatchString(stats.String()) {
		t.Fatalf("Wanted '%s'. Got '%s'.", want, stats)
	}
}

func TestCompareError(t *testing.T) {
	config := setupConfigError()
	want := "Searched 9 file(s), 4 file(s) different, 2 director(ies) different, 6 total entr(ies) different. 2 error(s)."

	stats, err := compare.Compare(config)
	fmt.Println(stats)

	if err != nil {
		t.Fatal("Error:", err)
	}
	if stats.String() != want {
		t.Fatalf("Wanted '%s'. Got '%s'.", want, stats)
	}
}

func BenchmarkCompareBasic(b *testing.B) {
	defer test.Quiet()()
	config := setupConfigBasic()

	for i := 0; i < b.N; i++ {
		compare.Compare(config)
	}
}

func BenchmarkCompareLarge(b *testing.B) {
	defer test.Quiet()()
	config := setupConfigLarge()

	for i := 0; i < b.N; i++ {
		compare.Compare(config)
	}
}

func BenchmarkCompareEntire(b *testing.B) {
	// Require the long flag to be set in order to run this long running benchmark.
	test.RequireLongBenchmark(b)

	defer test.Quiet()()
	config := setupConfigEntire()

	for i := 0; i < b.N; i++ {
		compare.Compare(config)
	}
}

var compareSameRegex = regexp.MustCompile(`Searched [[:digit:]]+ file.*, 0 file.* different, 0 director.* different, 0 total .* different. 0 error.*`)

func setupConfigBasic() *compare.Config {
	original := "../test/original"
	compared := "../test/compared"
	return compare.ParseConfigArgs("", []string{original, compared})
}

func setupConfigLarge() *compare.Config {
	large := "../test/large"
	return compare.ParseConfigArgs("", []string{large, large})
}

func setupConfigEntire() *compare.Config {
	large := "../test/large"
	return compare.ParseConfigArgs("", []string{"--entire", large, large})
}

func setupConfigError() *compare.Config {
	original := "../test/original"
	compared := "../test/errors"
	return compare.ParseConfigArgs("", []string{original, compared})
}
