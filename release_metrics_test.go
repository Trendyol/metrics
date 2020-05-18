package metrics

import (
	"testing"
	"time"
)

func TestReleaseMetrics_GetReleases(t *testing.T) {
	tests := []struct {
		name           string
		rs             ReleaseMetrics
		wantMetricsLen int
	}{
		{
			name: "Successful Releases",
			rs: ReleaseMetrics{
				{
					IsFix: false,
				},
				{
					IsFix: true,
				},
				{
					IsFix: false,
				},
			},
			wantMetricsLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rs.GetReleases(); len(got) != tt.wantMetricsLen {
				t.Errorf("GetReleases() = %v, want %v", len(got), tt.wantMetricsLen)
			}
		})
	}
}

func TestReleaseMetrics_GetFixReleases(t *testing.T) {
	tests := []struct {
		name           string
		rs             ReleaseMetrics
		wantMetricsLen int
	}{
		{
			name: "Fix Releases",
			rs: ReleaseMetrics{
				{
					IsFix: false,
				},
				{
					IsFix: true,
				},
				{
					IsFix: false,
				},
			},
			wantMetricsLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rs.GetFixReleases(); len(got) != tt.wantMetricsLen {
				t.Errorf("GetFixReleases() = %v, want %v", len(got), tt.wantMetricsLen)
			}
		})
	}
}

func TestReleaseMetrics_FilterByDate(t *testing.T) {
	type args struct {
		startDate time.Time
		endDate   time.Time
	}

	tests := []struct {
		name           string
		rs             ReleaseMetrics
		args           args
		wantMetricsLen int
	}{
		{
			name: "Before And After",
			rs: ReleaseMetrics{
				{
					ToDate: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				{
					ToDate: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				},
			},
			args: args{
				startDate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				endDate:   time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC),
			},
			wantMetricsLen: 2,
		},
		{
			name: "No Match",
			rs: ReleaseMetrics{
				{
					ToDate: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				{
					ToDate: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				},
			},
			args: args{
				startDate: time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC),
				endDate:   time.Date(2020, 1, 6, 0, 0, 0, 0, time.UTC),
			},
			wantMetricsLen: 0,
		},
		{
			name: "Before Equal And After Equal",
			rs: ReleaseMetrics{
				{
					ToDate: time.Date(2020, 1, 1, 11, 0, 0, 0, time.UTC),
				},
				{
					ToDate: time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
				},
			},
			args: args{
				startDate: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				endDate:   time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			wantMetricsLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rs.FilterByDate(tt.args.startDate, tt.args.endDate); len(got) != tt.wantMetricsLen {
				t.Errorf("FilterByDate() = %v, want %v", len(got), tt.wantMetricsLen)
			}
		})
	}
}

func TestReleaseMetrics_GetDeploymentFrequency(t *testing.T) {
	tests := []struct {
		name string
		rs   ReleaseMetrics
		want time.Duration
	}{
		{
			name: "Deployment Frequency",
			rs: ReleaseMetrics{
				{
					Interval: 1 * time.Hour,
				},
				{
					Interval: 2 * time.Hour,
				},
				{
					Interval: 3 * time.Hour,
				},
				{
					Interval: 4 * time.Hour,
				},
			},
			want: 150 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rs.GetDeploymentFrequency(); got != tt.want {
				t.Errorf("GetDeploymentFrequency() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReleaseMetrics_GetDeliveryLeadTime(t *testing.T) {
	tests := []struct {
		name string
		rs   ReleaseMetrics
		want time.Duration
	}{
		{
			name: "Delivery Lead Time",
			rs: ReleaseMetrics{
				{
					AverageCommitAge: 1 * time.Hour,
				},
				{
					AverageCommitAge: 2 * time.Hour,
				},
				{
					AverageCommitAge: 3 * time.Hour,
				},
				{
					AverageCommitAge: 4 * time.Hour,
				},
			},
			want: 150 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rs.GetDeliveryLeadTime(); got != tt.want {
				t.Errorf("GetDeliveryLeadTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReleaseMetrics_GetMeanTimeToRestore(t *testing.T) {
	tests := []struct {
		name string
		rs   ReleaseMetrics
		want time.Duration
	}{
		{
			name: "Mean Time To Restore",
			rs: ReleaseMetrics{
				{
					Interval: 1 * time.Hour,
					IsFix:    true,
				},
				{
					Interval: 2 * time.Hour,
					IsFix:    false,
				},
				{
					Interval: 3 * time.Hour,
					IsFix:    false,
				},
				{
					Interval: 4 * time.Hour,
					IsFix:    true,
				},
			},
			want: 150 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rs.GetMeanTimeToRestore(); got != tt.want {
				t.Errorf("GetMeanTimeToRestore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReleaseMetrics_GetChangeFailureRate(t *testing.T) {
	tests := []struct {
		name string
		rs   ReleaseMetrics
		want float64
	}{
		{
			name: "Change Failure Rate with no fix-releases",
			rs: ReleaseMetrics{
				{
					IsFix: false,
				},
				{
					IsFix: false,
				},
				{
					IsFix: false,
				},
				{
					IsFix: false,
				},
			},
			want: 0,
		},
		{
			name: "Change Failure Rate with fix-releases",
			rs: ReleaseMetrics{
				{
					IsFix: true,
				},
				{
					IsFix: false,
				},
				{
					IsFix: false,
				},
				{
					IsFix: false,
				},
			},
			want: 33.33,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rs.GetChangeFailureRate(); got != tt.want {
				t.Errorf("GetChangeFailRate() = %v, want %v", got, tt.want)
			}
		})
	}
}
