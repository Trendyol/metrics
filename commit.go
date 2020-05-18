package metrics

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Commits []Commit

type Commit struct {
	SHA  string
	Date time.Time
}

func (cs Commits) GetAverageCommitAge(start time.Time) time.Duration {
	if len(cs) == 0 {
		return time.Duration(0)
	}

	return start.Sub(cs[len(cs)/2].Date)
}

func ParseCommits(reader io.Reader) Commits {
	commits := make(Commits, 0)

	for scanner := bufio.NewScanner(reader); scanner.Scan(); {
		hashAndDate := strings.Split(scanner.Text(), ",")
		if len(hashAndDate) != 2 {
			continue
		}

		commitHash := hashAndDate[0]
		commitDate, err := time.Parse(time.RFC3339, hashAndDate[1]) // intentionally unhandled
		if err != nil {
			continue
		}

		commits = append(commits, Commit{
			SHA:  commitHash,
			Date: commitDate,
		})
	}

	return commits
}

func GetRawCommits(repoDir, startTag, endTag string) (io.Reader, error) {
	cmd := exec.Command("git", "log", `--pretty=format:%h,%aI`, fmt.Sprintf("%s..%s", startTag, endTag), `--no-merges`)
	cmd.Dir = repoDir
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git log failed. err: %w", err)
	}

	return bytes.NewReader(out), nil
}
