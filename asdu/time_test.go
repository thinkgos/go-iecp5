package asdu

import (
	"testing"
	"time"
)

var goldenCP24Time2as = []struct {
	serial         []byte
	min, sec, nsec int
	ok             bool
}{
	{[]byte{1, 2, 3}, 3, 0, 513000000, true},
	{[]byte{1, 2, 131}, 3, 0, 513000000, false},
	{[]byte{11, 12, 13}, 13, 3, 83000000, true},
}

func TestParseCP24Time2a(t *testing.T) {
	for _, gold := range goldenCP24Time2as {
		got := ParseCP24Time2a(gold.serial, time.UTC)

		switch {
		case !gold.ok && got != nil:
			t.Errorf("%#x: got %s for invalid", gold.serial, got)
		case !gold.ok:
			break
		case got == nil:
			t.Errorf("%#x: got nil for valid", gold.serial)
		case got.Nanosecond() != gold.nsec, got.Second() != gold.sec, got.Minute() != gold.min:
			t.Errorf("%#x: got %s, want …:%02d:%02d.%09d", gold.serial, got, gold.min, gold.sec, gold.nsec)

		}
	}
}

var goldenCP56Time2as = []struct {
	serial []byte
	time   time.Time
	ok     bool
}{
	{[]byte{1, 2, 3, 4, 5, 6, 7}, time.Date(2007, 6, 5, 4, 3, 0, 513000000, time.UTC), true},
	{[]byte{1, 2, 131, 4, 5, 6, 7}, time.Date(2007, 6, 5, 4, 3, 0, 513000000, time.UTC), false},
	{[]byte{11, 12, 13, 14, 15, 16, 17}, time.Date(2016, 12, 15, 14, 13, 3, 83000000, time.UTC), true},
}

func TestParseCP56Time2a(t *testing.T) {
	for _, gold := range goldenCP56Time2as {
		got := ParseCP56Time2a(gold.serial, time.UTC)

		switch {
		case !gold.ok && got != nil:
			t.Errorf("%#x: got %s for invalid", gold.serial, got)
		case gold.ok && got == nil:
			t.Errorf("%#x: got nil for valid", gold.serial)
		case gold.ok && *got != gold.time:
			t.Errorf("%#x: got %s, want %s", gold.serial, got, gold.time)
		}
	}
}
