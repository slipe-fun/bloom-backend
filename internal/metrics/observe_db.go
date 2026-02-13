package metrics

import "time"

func ObserveDB(query string, duration time.Duration, err error) {
	DbRequestsTotal.WithLabelValues(query).Inc()

	DbQueryDurationSeconds.
		WithLabelValues(query).
		Observe(duration.Seconds())

	if err != nil {
		DbErrorsTotal.WithLabelValues(query).Inc()
	}
}
