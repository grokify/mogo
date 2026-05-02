package duration

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDurationMillisecondsMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		duration DurationMilliseconds
		want     string
	}{
		{"zero", DurationMilliseconds(0), "0"},
		{"one_second", DurationMilliseconds(time.Second), "1000"},
		{"one_minute", DurationMilliseconds(time.Minute), "60000"},
		{"half_second", DurationMilliseconds(500 * time.Millisecond), "500"},
		{"complex", DurationMilliseconds(2*time.Minute + 30*time.Second + 500*time.Millisecond), "150500"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.duration.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("MarshalJSON() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestDurationMillisecondsUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    DurationMilliseconds
		wantErr bool
	}{
		{"zero", "0", DurationMilliseconds(0), false},
		{"one_second", "1000", DurationMilliseconds(time.Second), false},
		{"one_minute", "60000", DurationMilliseconds(time.Minute), false},
		{"with_spaces", " 500 ", DurationMilliseconds(500 * time.Millisecond), false},
		{"null", "null", DurationMilliseconds(0), false},
		{"empty", "", DurationMilliseconds(0), false},
		{"invalid", "abc", DurationMilliseconds(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d DurationMilliseconds
			err := d.UnmarshalJSON([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Fatalf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && d != tt.want {
				t.Errorf("UnmarshalJSON() = %v, want %v", d, tt.want)
			}
		})
	}
}

func TestDurationMillisecondsRoundTrip(t *testing.T) {
	type testStruct struct {
		Start DurationMilliseconds `json:"start_ms"`
		End   DurationMilliseconds `json:"end_ms"`
	}

	original := testStruct{
		Start: DurationMilliseconds(1500 * time.Millisecond),
		End:   DurationMilliseconds(3200 * time.Millisecond),
	}

	// Marshal
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	expectedJSON := `{"start_ms":1500,"end_ms":3200}`
	if string(data) != expectedJSON {
		t.Errorf("json.Marshal() = %s, want %s", data, expectedJSON)
	}

	// Unmarshal
	var decoded testStruct
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if decoded.Start != original.Start {
		t.Errorf("Start = %v, want %v", decoded.Start, original.Start)
	}
	if decoded.End != original.End {
		t.Errorf("End = %v, want %v", decoded.End, original.End)
	}
}

func TestDurationMillisecondsMethods(t *testing.T) {
	d := DurationMilliseconds(2500 * time.Millisecond)

	if got := d.Duration(); got != 2500*time.Millisecond {
		t.Errorf("Duration() = %v, want %v", got, 2500*time.Millisecond)
	}

	if got := d.Milliseconds(); got != 2500 {
		t.Errorf("Milliseconds() = %d, want %d", got, 2500)
	}

	if got := d.Seconds(); got != 2.5 {
		t.Errorf("Seconds() = %f, want %f", got, 2.5)
	}

	if got := d.String(); got != "2.5s" {
		t.Errorf("String() = %s, want %s", got, "2.5s")
	}
}

func TestFromDuration(t *testing.T) {
	d := time.Second + 500*time.Millisecond
	dm := FromDuration(d)
	if dm.Duration() != d {
		t.Errorf("FromDuration() = %v, want %v", dm.Duration(), d)
	}
}

func TestFromMilliseconds(t *testing.T) {
	dm := FromMilliseconds(1500)
	want := DurationMilliseconds(1500 * time.Millisecond)
	if dm != want {
		t.Errorf("FromMilliseconds() = %v, want %v", dm, want)
	}
}
