// Package histogram provides primitives for working with histograms,
// particularly towards building Prometheus exporters.
package histogram

import (
	"sync"

	"github.com/giantswarm/microerror"
)

type Config struct {
	// BucketLimits is the upper limit of each bucket in the histogram.
	// See https://godoc.org/github.com/prometheus/client_golang/prometheus#HistogramOpts.
	BucketLimits []float64
}

// Histogram is a data structure suitable for providing the inputs to a Prometheus Histogram.
// See https://godoc.org/github.com/prometheus/client_golang/prometheus#MustNewConstHistogram.
type Histogram struct {
	// count is the number of entries added to the Histogram.
	count uint32
	// sum is the total sum of all entries added to the Histogram.
	sum float64
	// buckets is a map of upper bounds to cumulative counts of entries,
	// excluding the +Inf bucket.
	buckets map[float64]uint32
	// mutex controls access to the buckets map, to allow for safe concurrent access.
	mutex sync.Mutex
}

func New(config Config) (*Histogram, error) {
	if len(config.BucketLimits) == 0 {
		return nil, microerror.Maskf(invalidConfigError, "%T.BucketLimits must not be empty", config)
	}

	buckets := map[float64]uint32{}
	for _, bucketLimit := range config.BucketLimits {
		buckets[bucketLimit] = 0
	}

	h := &Histogram{
		buckets: buckets,
	}

	return h, nil
}

// Add saves an entry to the Histogram.
func (h *Histogram) Add(x float64) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.count++
	h.sum += x

	for bucket, _ := range h.buckets {
		if x <= bucket {
			h.buckets[bucket]++
		}
	}
}

// Count returns the number of samples recorded.
func (h *Histogram) Count() uint64 {
	return uint64(h.count)
}

// Buckets returns the sum of all samples recorded.
func (h *Histogram) Sum() float64 {
	return h.sum
}

// Buckets returns a copy of the current buckets with their counts.
func (h *Histogram) Buckets() map[float64]uint64 {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	bucketsCopy := map[float64]uint64{}

	for value, count := range h.buckets {
		bucketsCopy[value] = uint64(count)
	}

	return bucketsCopy
}
