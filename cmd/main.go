package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"metrics"
)

var (
	releasesTagRegexPattern = flag.String("r", "releases/\\d+", "")
	fixTagRegexPattern      = flag.String("f", "releases/fix/\\d+", "")
	startDate               = flag.String("s", "", "")
	sizeDays                = flag.Int("d", 7, "")
)

var usage = `Usage: metrics [options...] [repositories...]

Options:
  -r Releases tag regex pattern. Default is: "releases/\\d+"
  -f Fix tag regex pattern. Default is: "releases/fix/\\d+"
  -s Earliest start date of sprint. Date format must be RFC3339. For example: 2020-04-22T16:00:00+03:00
  -d Sprint days size. Default is 7 days

For Example: metrics -s 2020-04-22T16:00:00+03:00 ~/code/my-api ~/code/my-api-2
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	flag.Parse()

	if flag.NArg() < 1 {
		usageAndExit()
	}

	sprintStartDate, err := time.Parse(time.RFC3339, *startDate)
	if err != nil {
		errAndExit(err.Error())
	}

	repositories := flag.Args()

	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	records := [][]string{
		{"repository", "from_date", "to_date", "deployment_frequency", "delivery_lead_time", "mean_time_to_restore", "change_failure_rate"},
	}

	for _, repository := range repositories {
		releases, err := metrics.GetReleases(repository, *releasesTagRegexPattern, *fixTagRegexPattern)
		if err != nil {
			errAndExit(err.Error())
		}

		releaseMetrics := metrics.GetReleaseMetrics(releases).GetFourKeyMetricsForSprints(sprintStartDate, *sizeDays)

		for _, releaseMetric := range releaseMetrics {
			records = append(records, []string{
				path.Base(repository),
				releaseMetric.ToDate.String(),
				releaseMetric.FromDate.String(),
				releaseMetric.DeploymentFrequency.String(),
				releaseMetric.DeliveryLeadTime.String(),
				releaseMetric.MeanTimeToRestore.String(),
				fmt.Sprintf("%.2f", releaseMetric.ChangeFailureRate),
			})
		}
	}

	for _, record := range records {
		if err := w.Write(record); err != nil {
			errAndExit(err.Error())
		}
	}
}

func errAndExit(msg string) {
	fmt.Fprintf(os.Stderr, msg)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func usageAndExit() {
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}
