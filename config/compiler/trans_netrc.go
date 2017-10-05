package compiler

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/config"
	"github.com/drone/drone-yaml-v1/yaml"
	"github.com/vincent-petithory/dataurl"
)

const (
	netrcPath = "/root/netrc"
	netrcMode = 0600
)

func transformNetrc(machine, username, password string) Transform {
	return func(dst *engine.Step, src *yaml.Container, conf *config.Config) {
		if machine+username+password == "" {
			return
		}
		netrc := generateNetrc(machine, username, password)
		tarball := generateTarball(netrcPath, netrc, netrcMode)
		datauri := dataurl.EncodeBytes(tarball)
		dst.Restore = append(dst.Restore, &engine.Snapshot{
			Data:   []byte(datauri),
			Target: "/",
		})
	}
}

func generateTarball(filepath, filedata string, filemode int64) []byte {
	b := new(bytes.Buffer)
	t := tar.NewWriter(b)
	h := &tar.Header{
		Name:    filepath,
		Mode:    filemode,
		Size:    int64(len(filedata)),
		ModTime: time.Now(),
	}
	t.WriteHeader(h)
	io.WriteString(t, filedata)
	return b.Bytes()
}

func generateNetrc(machine, username, password string) string {
	return fmt.Sprintf("machine %s login %s password %s",
		machine, username, password)
}
