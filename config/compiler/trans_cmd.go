package compiler

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/config"
	"github.com/drone/drone-yaml-v1/yaml"

	"github.com/vincent-petithory/dataurl"
)

func transformCommand(dst *engine.Step, src *yaml.Container, conf *config.Config) {
	if conf.Platform.Name == "windows/amd64" {
		transformCommandWin(dst, src, conf)
		return
	}
	if len(src.Commands) == 0 {
		return
	}
	shell := src.Shell
	if len(src.Shell) == 0 {
		shell = "/bin/sh"
	}

	script := generateScriptPosix(src.Commands)
	dst.Entrypoint = []string{shell}
	dst.Command = []string{"/bin/_drone"}
	dst.Restore = append(dst.Restore, &engine.Snapshot{
		Source: script,
		Target: "/",
	})
}

// generateScriptPosix is a helper function that generates a build script
// for a linux container using the given
func generateScriptPosix(commands []string) string {
	var buf bytes.Buffer
	for _, command := range commands {
		escaped := fmt.Sprintf("%q", command)
		escaped = strings.Replace(escaped, "$", `\$`, -1)
		buf.WriteString(fmt.Sprintf(
			traceScript,
			escaped,
			command,
		))
	}
	script := fmt.Sprintf(
		setupScript,
		buf.String(),
	)

	w := new(bytes.Buffer)
	t := tar.NewWriter(w)
	h := &tar.Header{
		Name:    "bin/_drone",
		Mode:    0644,
		Size:    int64(len(script)),
		ModTime: time.Now(),
	}
	t.WriteHeader(h)
	io.WriteString(t, script)
	return dataurl.EncodeBytes(w.Bytes()) // .New(w.Bytes(), "application/x-tar", nil).String()
}

// setupScript is a helper script this is added to the build to ensure
// a minimum set of environment variables are set correctly.
const setupScript = `
if [ -n "$CI_NETRC_MACHINE" ]; then
cat <<EOF > $HOME/.netrc
machine $CI_NETRC_MACHINE
login $CI_NETRC_USERNAME
password $CI_NETRC_PASSWORD
EOF
chmod 0600 $HOME/.netrc
fi
unset CI_NETRC_USERNAME
unset CI_NETRC_PASSWORD
unset CI_SCRIPT
unset DRONE_NETRC_USERNAME
unset DRONE_NETRC_PASSWORD

set -e

%s
`

// traceScript is a helper script that is added to the build script
// to trace a command.
const traceScript = `
echo + %s
%s
`
