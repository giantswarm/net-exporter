// Package histogramvec provides primitives for working with a vector of histograms,
// particularly towards building Prometheus exporters.
package histogramvec

import (
	"sync"

	"github.com/giantswarm/exporterkit/histogram"
	"github.com/giantswarm/microerror"
)

type Config struct {
	// BucketLimits is the upper limit of each bucket used when creating new internal Histograms.
	// See https://godoc.org/github.com/prometheus/client_golang/prometheus#HistogramOpts.
	BucketLimits []float64
}

// HistogramVec is a data structure suitable for holding multiple Histograms,
// and then providing the inputs to multiple Prometheus Histograms.
// See https://godoc.org/github.com/prometheus/client_golang/prometheus#MustNewConstHistogram.
type HistogramVec struct {
	// bucketLimits is the upper limit of buckets for the Histograms that the HistogramVec manages.
	bucketLimits []float64

	// histograms is a mapping between labels and the Histograms that the HistogramVec manages.
	histograms map[string]*histogram.Histogram
	// mutex controls access to the histograms map, to allow for safe concurrent access.
	mutex sync.Mutex
}

func New(config Config) (*HistogramVec, error) {
	if len(config.BucketLimits) == 0 {
		return nil, microerror.Maskf(invalidConfigError, "%T.BucketLimits must not be empty", config)
	}

	hv := &HistogramVec{
		bucketLimits: config.BucketLimits,

		histograms: map[string]*histogram.Histogram{},
	}

	return hv, nil
}

// Add saves an entry to the Histogram with the given label,
// creating it internally if required.
func (hv *HistogramVec) Add(label string, x float64) error {
	hv.mutex.Lock()
	defer hv.mutex.Unlock()

	if _, ok := hv.histograms[label]; !ok {
		c := histogram.Config{
			BucketLimits: hv.bucketLimits,
		}

		h, err := histogram.New(c)
		if err != nil {
			return microerror.Mask(err)
		}

		hv.histograms[label] = h
	}

	hv.histograms[label].Add(x)

	return nil
}

// Ensure removes any internal Histograms that aren't in the given slice of labels.
// This is useful when a label is no longer being recorded, such as in dynamic systems.
func (hv *HistogramVec) Ensure(labels []string) {
	hv.mutex.Lock()
	defer hv.mutex.Unlock()

	for existingLabel := range hv.histograms {
		labelRequested := false

		for _, requestedLabel := range labels {
			if requestedLabel == existingLabel {
				labelRequested = true
			}
		}

		if !labelRequested {
			delete(hv.histograms, existingLabel)
		}
	}
}

// Histograms returns a copy of the currently managed Histograms.
func (hv *HistogramVec) Histograms() map[string]*histogram.Histogram {
	hv.mutex.Lock()
	defer hv.mutex.Unlock()

	histogramsCopy := map[string]*histogram.Histogram{}

	for label, histogram := range hv.histograms {
		h := *histogram
		histogramsCopy[label] = &h
	}

	return histogramsCopy
}
