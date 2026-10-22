package network

import (
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	ip1 = "127.0.0.1"
	ip2 = "127.0.0.2"
	ip3 = "127.0.0.3"
)

func Test_calculateNeighbours(t *testing.T) {
	testCases := []struct {
		name               string
		inputN             int
		inputIP            string
		inputAddresses     []string
		expectedNeighbours []string
	}{
		{
			name:               "case 0",
			inputN:             1,
			inputIP:            ip1,
			inputAddresses:     []string{},
			expectedNeighbours: []string{},
		},
		{
			name:               "case 1",
			inputN:             1,
			inputIP:            ip1,
			inputAddresses:     []string{ip1},
			expectedNeighbours: []string{ip1},
		},
		{
			name:               "case 2",
			inputN:             1,
			inputIP:            ip1,
			inputAddresses:     []string{ip1, ip2},
			expectedNeighbours: []string{ip2},
		},
		{
			name:               "case 3",
			inputN:             1,
			inputIP:            ip2,
			inputAddresses:     []string{ip1, ip2},
			expectedNeighbours: []string{ip1},
		},
		{
			name:               "case 4",
			inputN:             2,
			inputIP:            ip1,
			inputAddresses:     []string{},
			expectedNeighbours: []string{},
		},
		{
			name:               "case 5",
			inputN:             2,
			inputIP:            ip1,
			inputAddresses:     []string{ip1},
			expectedNeighbours: []string{ip1},
		},
		{
			name:               "case 6",
			inputN:             2,
			inputIP:            ip1,
			inputAddresses:     []string{ip1, ip2},
			expectedNeighbours: []string{ip2, ip1},
		},
		{
			name:               "case 7",
			inputN:             2,
			inputIP:            ip2,
			inputAddresses:     []string{ip1, ip2, ip3},
			expectedNeighbours: []string{ip3, ip1},
		},
		{
			name:               "case 8",
			inputN:             2,
			inputIP:            "127.0.0.4",
			inputAddresses:     []string{ip1, ip2, ip3, "127.0.0.4"},
			expectedNeighbours: []string{ip1, ip2},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			collector := Collector{}
			returnedNeighbours := collector.calculateNeighbours(tc.inputN, tc.inputIP, tc.inputAddresses)

			if !cmp.Equal(returnedNeighbours, tc.expectedNeighbours) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.expectedNeighbours, returnedNeighbours))
			}
		})
	}
}
