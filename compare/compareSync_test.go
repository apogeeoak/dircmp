package compare_test

import (
	"fmt"
	"testing"

	"github.com/apogeeoak/dircmp/compare"
	"github.com/apogeeoak/dircmp/lib/test"
)

func TestCompareSyncBasic(t *testing.T) {
	config := setupConfigBasic()
	want := "Searched 8 file(s), 4 file(s) different, 2 director(ies) different, 6 total entr(ies) different."

	stats, err := compare.CompareSync(config)
	fmt.Println(stats)

	if err != nil {
		t.Fatal("Error:", err)
	}
	if stats.String() != want {
		t.Fatalf("Wanted '%s'. Got '%s'.", want, stats)
	}
}

func TestCompareSyncLarge(t *testing.T) {
	config := setupConfigLarge()
	want := compareSameRegex

	stats, err := compare.CompareSync(config)
	fmt.Println(stats)

	if err != nil {
		t.Fatal("Error:", err)
	}
	if !want.MatchString(stats.String()) {
		t.Fatalf("Wanted '%s'. Got '%s'.", want, stats)
	}
}

func TestCompareSyncEntire(t *testing.T) {
	// Require the long flag to be set in order to run this long running test.
	test.RequireLong(t)

	config := setupConfigEntire()
	want := compareSameRegex

	stats, err := compare.CompareSync(config)
	fmt.Println(stats)

	if err != nil {
		t.Fatal("Error:", err)
	}
	if !want.MatchString(stats.String()) {
		t.Fatalf("Wanted '%s'. Got '%s'.", want, stats)
	}
}

func BenchmarkCompareSyncBasic(b *testing.B) {
	defer test.Quiet()()
	config := setupConfigBasic()

	for i := 0; i < b.N; i++ {
		compare.CompareSync(config)
	}
}

func BenchmarkCompareSyncLarge(b *testing.B) {
	defer test.Quiet()()
	config := setupConfigLarge()

	for i := 0; i < b.N; i++ {
		compare.CompareSync(config)
	}
}

func BenchmarkCompareSyncEntire(b *testing.B) {
	// Require the long flag to be set in order to run this long running benchmark.
	test.RequireLongBenchmark(b)

	defer test.Quiet()()
	config := setupConfigEntire()

	for i := 0; i < b.N; i++ {
		compare.CompareSync(config)
	}
}
