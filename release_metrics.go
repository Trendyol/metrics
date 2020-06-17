package metrics

import (
	"math"
	"time"
)

type ReleaseMetrics []ReleaseMetric

type ReleaseMetric struct {
	From             string
	To               string
	FromDate         time.Time
	ToDate           time.Time
	Interval         time.Duration
	IsFix            bool
	AverageCommitAge time.Duration
}

func GetReleaseMetrics(releases []Release) ReleaseMetrics {
	releaseMetrics := make(ReleaseMetrics, 0)

	for i, j := 0, 1; j < len(releases); i, j = i+1, j+1 {
		currentDeployment, nextDeployment := releases[i], releases[j]
		releaseMetrics = append(releaseMetrics, ReleaseMetric{
			From:             nextDeployment.Tag,
			To:               currentDeployment.Tag,
			FromDate:         nextDeployment.Date,
			ToDate:           currentDeployment.Date,
			Interval:         currentDeployment.Date.Sub(nextDeployment.Date),
			IsFix:            currentDeployment.IsFix,
			AverageCommitAge: currentDeployment.Commits.GetAverageCommitAge(currentDeployment.Date),
		})
	}

	return releaseMetrics
}

func (rs ReleaseMetrics) GetReleases() ReleaseMetrics {
	releaseMetrics := make(ReleaseMetrics, 0)

	for _, r := range rs {
		if !r.IsFix {
			releaseMetrics = append(releaseMetrics, r)
		}
	}

	return releaseMetrics
}

func (rs ReleaseMetrics) GetFixReleases() ReleaseMetrics {
	fixReleaseMetrics := make(ReleaseMetrics, 0)

	for _, r := range rs {
		if r.IsFix {
			fixReleaseMetrics = append(fixReleaseMetrics, r)
		}
	}

	return fixReleaseMetrics
}

func (rs ReleaseMetrics) FilterByDate(startDate time.Time, endDate time.Time) ReleaseMetrics {
	filteredMetrics := make(ReleaseMetrics, 0)

	for _, r := range rs {
		if (r.ToDate.After(startDate) || r.ToDate.Equal(startDate)) && (r.ToDate.Before(endDate) || r.ToDate.Equal(endDate)) {
			filteredMetrics = append(filteredMetrics, r)
		}
	}

	return filteredMetrics
}

func (rs ReleaseMetrics) GetDeploymentFrequency() time.Duration {
	var totalDeploymentInterval time.Duration
	for _, r := range rs {
		totalDeploymentInterval += r.Interval
	}

	avgDeploymentInterval := totalDeploymentInterval.Seconds() / float64(len(rs))
	return time.Second * time.Duration(avgDeploymentInterval)
}

func (rs ReleaseMetrics) GetDeliveryLeadTime() time.Duration {
	var totalAverageCommitAge time.Duration
	for _, r := range rs {
		totalAverageCommitAge += r.AverageCommitAge
	}

	avgCommitAge := totalAverageCommitAge.Seconds() / float64(len(rs))
	return time.Second * time.Duration(avgCommitAge)
}

func (rs ReleaseMetrics) GetMeanTimeToRestore() time.Duration {
	fixReleases := rs.GetFixReleases()

	var totalFixReleasesDuration time.Duration
	for _, r := range fixReleases {
		totalFixReleasesDuration += r.Interval
	}

	avgFixReleasesDuration := totalFixReleasesDuration.Seconds() / float64(len(fixReleases))
	return time.Second * time.Duration(avgFixReleasesDuration)
}

func (rs ReleaseMetrics) GetChangeFailureRate() float64 {
	failRate := float64(len(rs.GetFixReleases())) / float64(len(rs.GetReleases())) * 100

	roundedFailRate := math.Trunc(failRate*100) / 100
	if math.IsNaN(roundedFailRate) {
		return 0
	}

	return roundedFailRate
}

type FourKeyMetrics struct {
	FromDate            time.Time
	ToDate              time.Time
	DeploymentFrequency time.Duration
	DeliveryLeadTime    time.Duration
	MeanTimeToRestore   time.Duration
	ChangeFailureRate   float64
}

func (rs ReleaseMetrics) GetFourKeyMetricsForSprints(earliestDate time.Time, sprintSizeDays int) (metrics []FourKeyMetrics) {
	for startDate := earliestDate; startDate.Before(time.Now()); startDate = startDate.AddDate(0, 0, sprintSizeDays) {
		endDate := startDate.AddDate(0, 0, sprintSizeDays)
		releaseMetrics := rs.FilterByDate(startDate, endDate)

		metrics = append(metrics, FourKeyMetrics{
			FromDate:            startDate,
			ToDate:              endDate,
			DeploymentFrequency: releaseMetrics.GetDeploymentFrequency(),
			DeliveryLeadTime:    releaseMetrics.GetDeliveryLeadTime(),
			MeanTimeToRestore:   releaseMetrics.GetMeanTimeToRestore(),
			ChangeFailureRate:   releaseMetrics.GetChangeFailureRate(),
		})
	}

	return
}

func (rs ReleaseMetrics) GetFourKeyMetricsLookBack(startDate time.Time, earliestDays int) FourKeyMetrics {
	earliestDate := startDate.AddDate(0, 0, -earliestDays)

	releaseMetrics := rs.FilterByDate(earliestDate, startDate)

	return FourKeyMetrics{
		FromDate:            earliestDate,
		ToDate:              startDate,
		DeploymentFrequency: releaseMetrics.GetDeploymentFrequency(),
		DeliveryLeadTime:    releaseMetrics.GetDeliveryLeadTime(),
		MeanTimeToRestore:   releaseMetrics.GetMeanTimeToRestore(),
		ChangeFailureRate:   releaseMetrics.GetChangeFailureRate(),
	}
}
