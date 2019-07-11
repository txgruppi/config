package config

import (
	"os"
	"path"
	"path/filepath"
)

func getConfigDir() (string, error) {
	dir := os.Getenv("CONFIG_DIR")
	if dir == "" {
		dir = "./config"
	}
	return filepath.Abs(dir)
}

func getDeployment() (string, error) {
	return os.Getenv("ENV"), nil
}

func makePossibleFilepaths(configDir, hostname, deployment string, exts []string) []string {
	out := []string{}
	nms := []string{"default"}
	if deployment != "" {
		nms = append(nms, deployment)
	}
	if hostname != "" {
		nms = append(nms, hostname)
	}
	if hostname != "" && deployment != "" {
		nms = append(nms, hostname+"-"+deployment)
	}
	nms = append(nms, "local")
	if deployment != "" {
		nms = append(nms, "local-"+deployment)
	}
	for _, nm := range nms {
		for _, ext := range exts {
			out = append(out, path.Join(configDir, nm+ext))
		}
	}

	return out
}
