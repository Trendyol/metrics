package metrics

import (
	"io"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParseTags(t *testing.T) {
	tests := []struct {
		name               string
		reader             io.Reader
		releasesTagPattern string
		fixTagRegex        string
		wantReleases       []Release
	}{
		{
			name:               "Tags",
			reader:             strings.NewReader("bla-bla,2020-05-06T11:58:40Z\nreleases/393,2020-05-06T11:58:40Z\nreleases/fix/374,2020-04-29T15:40:38Z"),
			releasesTagPattern: "releases/\\d+",
			fixTagRegex:        "releases/fix/\\d+",
			wantReleases: []Release{
				{
					Tag:   "releases/393",
					Date:  time.Date(2020, 5, 6, 11, 58, 40, 0, time.UTC),
					IsFix: false,
				},
				{
					Tag:   "releases/fix/374",
					Date:  time.Date(2020, 4, 29, 15, 40, 38, 0, time.UTC),
					IsFix: true,
				},
			},
		},
		{
			name:               "No Tags",
			reader:             strings.NewReader(""),
			releasesTagPattern: "releases/\\d+",
			fixTagRegex:        "releases/fix/\\d+",
			wantReleases:       make([]Release, 0),
		},
		{
			name:               "Malformed Raw Log",
			reader:             strings.NewReader("bla-bla-bla"),
			releasesTagPattern: "releases/\\d+",
			fixTagRegex:        "releases/fix/\\d+",
			wantReleases:       make([]Release, 0),
		},
		{
			name:               "Malformed Raw Date",
			reader:             strings.NewReader("releases/393,bla-bla-bla"),
			releasesTagPattern: "releases/\\d+",
			fixTagRegex:        "releases/fix/\\d+",
			wantReleases:       make([]Release, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotReleases := ParseReleases(tt.reader, tt.releasesTagPattern, tt.fixTagRegex); !reflect.DeepEqual(gotReleases, tt.wantReleases) {
				t.Errorf("ParseTags() gotReleases = %v, want %v", gotReleases, tt.wantReleases)
			}
		})
	}
}
