package metrics

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type Releases []Release

type Release struct {
	Tag     string
	IsFix   bool
	Date    time.Time
	Commits Commits
}

func GetReleases(repoDir, releasesTagRegex, fixTagRegex string) ([]Release, error) {
	tags, err := GetRawReleases(repoDir)
	if err != nil {
		return nil, err
	}

	releases := ParseReleases(tags, releasesTagRegex, fixTagRegex)

	for current, next := 0, 1; next < len(releases); current, next = current+1, next+1 {
		reader, err := GetRawCommits(repoDir, releases[next].Tag, releases[current].Tag)
		if err != nil {
			return nil, fmt.Errorf("couldn't get commits for %s..%s", releases[next].Tag, releases[current].Tag)
		}

		releases[current].Commits = ParseCommits(reader)
	}

	return releases, nil
}

func ParseReleases(reader io.Reader, releasesTagPattern, fixTagRegex string) []Release {
	releases := make([]Release, 0)
	releasesRegexp, fixRegexp := regexp.MustCompile(releasesTagPattern), regexp.MustCompile(fixTagRegex)

	for scanner := bufio.NewScanner(reader); scanner.Scan(); {
		tagAndDate := strings.Split(scanner.Text(), ",")
		if len(tagAndDate) != 2 {
			continue
		}

		tag := tagAndDate[0]
		tagDate, err := time.Parse(time.RFC3339, tagAndDate[1])
		if err != nil {
			continue
		}

		isReleases, isFix := releasesRegexp.MatchString(tag), fixRegexp.MatchString(tag)
		if !isReleases && !isFix {
			continue
		}

		releases = append(releases, Release{
			Tag:   tag,
			IsFix: isFix,
			Date:  tagDate,
		})
	}

	return releases
}

func GetRawReleases(repoDir string) (io.Reader, error) {
	cmd := exec.Command("git", "for-each-ref", "--sort=-creatordate", `--format=%(refname),%(creatordate:iso-strict)`, `refs/tags`)
	cmd.Stderr = os.Stderr
	cmd.Dir = repoDir

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git log failed. err: %w", err)
	}

	return bytes.NewReader(out), nil
}
