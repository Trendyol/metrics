package metrics

import (
	"io"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestCommits_GetAverageCommitAge(t *testing.T) {
	tests := []struct {
		name    string
		commits Commits
		start   time.Time
		want    time.Duration
	}{
		{
			name: "Calculate with Three Different Dates",
			commits: Commits{
				{
					Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					Date: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				{
					Date: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				},
			},
			start: time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC),
			want:  48 * time.Hour,
		},
		{
			name: "Calculate with Three Different Dates and Equal Date",
			commits: Commits{
				{
					Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					Date: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				{
					Date: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				},
			},
			start: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
			want:  24 * time.Hour,
		},
		{
			name:    "Calculate with No Commits",
			commits: make(Commits, 0),
			start:   time.Date(2020, 1, 1, 13, 0, 0, 0, time.UTC),
			want:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.commits.GetAverageCommitAge(tt.start); got != tt.want {
				t.Errorf("GetAverageCommitAge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseCommits(t *testing.T) {
	tests := []struct {
		name        string
		reader      io.Reader
		wantCommits Commits
	}{
		{
			name:   "Commits",
			reader: strings.NewReader("1,2020-04-29T15:30:00Z\n2,2020-05-29T15:30:00Z"),
			wantCommits: Commits{
				{
					SHA:  "1",
					Date: time.Date(2020, 4, 29, 15, 30, 0, 0, time.UTC),
				},
				{
					SHA:  "2",
					Date: time.Date(2020, 5, 29, 15, 30, 0, 0, time.UTC),
				},
			},
		},
		{
			name:        "No Commits",
			reader:      strings.NewReader(""),
			wantCommits: make(Commits, 0),
		},
		{
			name:        "Malformed Raw Log",
			reader:      strings.NewReader("bla-bla-bla"),
			wantCommits: make(Commits, 0),
		},
		{
			name:        "Malformed Date",
			reader:      strings.NewReader("1,any"),
			wantCommits: make(Commits, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if gotCommits := ParseCommits(tt.reader); !reflect.DeepEqual(gotCommits, tt.wantCommits) {
				t.Errorf("ParseCommits() = %v, want %v", gotCommits, tt.wantCommits)
			}
		})
	}
}
