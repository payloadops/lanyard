package util

import "fmt"

func StageRegionDisambiguator(prefix string, stage string, region string) string {
	return fmt.Sprintf("%s-%s-%s", prefix, stage, region)
}
