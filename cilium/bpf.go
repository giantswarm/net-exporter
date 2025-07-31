package cilium

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cilium/cilium/pkg/bpf"
	"github.com/cilium/cilium/pkg/maps/policymap"

	"github.com/giantswarm/microerror"
)

type policyMap struct {
	EndpointID string
	Path       string
	Content    policymap.PolicyEntriesDump
	Size       int
}

func mapContent(file string) (policymap.PolicyEntriesDump, error) {
	m, err := policymap.OpenPolicyMap(nil, file)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	defer m.Close()

	statsMap, err := m.DumpToSlice()
	if err != nil {
		return nil, microerror.Mask(err)
	}
	sort.Slice(statsMap, statsMap.Less)

	return statsMap, nil
}

func listAllMaps() ([]policyMap, error) {
	mapRootPrefixPath := bpf.TCGlobalsPath()
	mapMatchExpr := filepath.Join(mapRootPrefixPath, "cilium_policy_*")

	matchFiles, err := filepath.Glob(mapMatchExpr)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	if len(matchFiles) == 0 {
		fmt.Println("no maps found")
		return nil, nil
	}

	maps := []policyMap{}
	for _, file := range matchFiles {
		endpointSplit := strings.Split(file, "_")
		endpoint := strings.TrimLeft(endpointSplit[len(endpointSplit)-1], "0")
		mcontent, err := mapContent(file)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		maps = append(maps, policyMap{
			EndpointID: endpoint,
			Path:       file,
			Content:    mcontent,
			Size:       len(mcontent),
		})
	}

	return maps, nil
}
