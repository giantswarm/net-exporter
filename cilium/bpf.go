package cilium

import (
	"log/slog"
	"os"
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

func (c *Collector) mapContent(logger *slog.Logger, file string) (policymap.PolicyEntriesDump, error) {
	c.logger.Log("level", "info", "message", "opening policy map", "file", file)
	m, err := policymap.OpenPolicyMap(logger, file)
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

func (c *Collector) listAllMaps() ([]policyMap, error) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	mapRootPrefixPath := bpf.TCGlobalsPath()
	mapMatchExpr := filepath.Join(mapRootPrefixPath, "cilium_policy_*")

	matchFiles, err := filepath.Glob(mapMatchExpr)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	if len(matchFiles) == 0 {
		c.logger.Log("level", "info", "message", "no maps found", "path", mapMatchExpr)
		return nil, nil
	}

	maps := []policyMap{}
	for _, file := range matchFiles {
		endpointSplit := strings.Split(file, "_")
		endpoint := strings.TrimLeft(endpointSplit[len(endpointSplit)-1], "0")
		mcontent, err := c.mapContent(log, file)
		if err != nil {
			c.logger.Log("level", "info", "message", "no map found", "path", file, "error", err)
			continue
		}
		maps = append(maps, policyMap{
			EndpointID: endpoint,
			Path:       file,
			Content:    mcontent,
			Size:       len(mcontent),
		})
		c.logger.Log("level", "info", "message", "processed policy map file", "file", file, "endpoint", endpoint, "size", len(mcontent))
	}

	return maps, nil
}
