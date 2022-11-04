package clock

import (
	"os"
	"reflect"
	"testing"
	"time"
)

var testFake Clock

func TestMain(m *testing.M) {
	testFake = NewFakeClock(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC))

	code := m.Run()

	os.Exit(code)

}

func TestFakeClock_Now(t *testing.T) {
	tests := []struct {
		name string
		want time.Time
	}{
		{"normal", time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testFake.Now(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Now() = %v, want %v", got, tt.want)
			}
		})
	}
}
