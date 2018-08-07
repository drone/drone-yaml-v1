package compiler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/config"
	"github.com/drone/drone-yaml-v1/yaml"
)

func transformCommandWin(dst *engine.Step, src *yaml.Container, conf *config.Config) {
	if len(src.Commands) == 0 {
		return
	}
	if dst.Environment == nil {
		dst.Environment = map[string]string{}
	}

	dst.Entrypoint = []string{"powershell", "-noprofile", "-noninteractive", "-command"}
	dst.Command = []string{"[System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String($Env:CI_SCRIPT)) | iex"}
	dst.Environment["CI_SCRIPT"] = generateScriptWindows(src.Commands)
	dst.Environment["HOME"] = "c:\\root"
	dst.Environment["SHELL"] = "powershell.exe"
}

func generateScriptWindows(commands []string) string {
	var buf bytes.Buffer
	for _, command := range commands {
		escaped := fmt.Sprintf("%q", command)
		escaped = strings.Replace(escaped, "$", `\$`, -1)
		buf.WriteString(fmt.Sprintf(
			traceScriptWin,
			escaped,
			command,
		))
	}
	script := fmt.Sprintf(
		setupScriptWin,
		buf.String(),
	)
	return base64.StdEncoding.EncodeToString([]byte(script))
}

const setupScriptWin = `
$ErrorActionPreference = 'Stop';
&cmd /c "mkdir c:\root";
if ($Env:CI_NETRC_MACHINE) {
$netrc=[string]::Format("{0}\_netrc",$Env:HOME);
"machine $Env:CI_NETRC_MACHINE" >> $netrc;
"login $Env:CI_NETRC_USERNAME" >> $netrc;
"password $Env:CI_NETRC_PASSWORD" >> $netrc;
};
[Environment]::SetEnvironmentVariable("CI_NETRC_PASSWORD",$null);
[Environment]::SetEnvironmentVariable("CI_SCRIPT",$null);
[Environment]::SetEnvironmentVariable("DRONE_NETRC_USERNAME",$null);
[Environment]::SetEnvironmentVariable("DRONE_NETRC_PASSWORD",$null);
%s
`

// traceScript is a helper script that is added to the build script
// to trace a command.
const traceScriptWin = `
Write-Output ('+ %s');
& %s; if ($LASTEXITCODE -ne 0) {exit $LASTEXITCODE}
`
