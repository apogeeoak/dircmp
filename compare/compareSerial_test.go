package compare_test

import (
	"fmt"
	"testing"

	"github.com/apogeeoak/dircmp/compare"
	"github.com/apogeeoak/dircmp/lib/test"
)

func TestCompareSerialBasic(t *testing.T) {
	config := setupConfigBasic()
	want := "Searched 8 file(s), 4 file(s) different, 2 director(ies) different, 6 total entr(ies) different. 0 error(s)."

	stats, err := compare.CompareSerial(config)
	fmt.Println(stats)

	if err != nil {
		t.Fatal("Error:", err)
	}
	if stats.String() != want {
		t.Fatalf("Wanted '%s'. Got '%s'.", want, stats)
	}
}

func TestCompareSerialLarge(t *testing.T) {
	config := setupConfigLarge()
	want := compareSameRegex

	stats, err := compare.CompareSerial(config)
	fmt.Println(stats)

	if err != nil {
		t.Fatal("Error:", err)
	}
	if !want.MatchString(stats.String()) {
		t.Fatalf("Wanted '%s'. Got '%s'.", want, stats)
	}
}

func TestCompareSerialEntire(t *testing.T) {
	// Require the long flag to be set in order to run this long running test.
	test.RequireLong(t)

	config := setupConfigEntire()
	want := compareSameRegex

	stats, err := compare.CompareSerial(config)
	fmt.Println(stats)

	if err != nil {
		t.Fatal("Error:", err)
	}
	if !want.MatchString(stats.String()) {
		t.Fatalf("Wanted '%s'. Got '%s'.", want, stats)
	}
}

func TestCompareSerialError(t *testing.T) {
	config := setupConfigError()
	want := "Searched 9 file(s), 4 file(s) different, 2 director(ies) different, 6 total entr(ies) different. 2 error(s)."

	stats, err := compare.CompareSerial(config)
	fmt.Println(stats)

	if err != nil {
		t.Fatal("Error:", err)
	}
	if stats.String() != want {
		t.Fatalf("Wanted '%s'. Got '%s'.", want, stats)
	}
}

func BenchmarkCompareSerialBasic(b *testing.B) {
	defer test.Quiet()()
	config := setupConfigBasic()

	for i := 0; i < b.N; i++ {
		compare.CompareSerial(config)
	}
}

func BenchmarkCompareSerialLarge(b *testing.B) {
	defer test.Quiet()()
	config := setupConfigLarge()

	for i := 0; i < b.N; i++ {
		compare.CompareSerial(config)
	}
}

func BenchmarkCompareSerialEntire(b *testing.B) {
	// Require the long flag to be set in order to run this long running benchmark.
	test.RequireLongBenchmark(b)

	defer test.Quiet()()
	config := setupConfigEntire()

	for i := 0; i < b.N; i++ {
		compare.CompareSerial(config)
	}
}
