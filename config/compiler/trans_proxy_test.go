package compiler

import (
	"os"
	"testing"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/yaml"
)

func Test_transformProxy(t *testing.T) {
	var (
		noProxy    = getenv("no_proxy")
		httpProxy  = getenv("https_proxy")
		httpsProxy = getenv("https_proxy")
	)
	defer func() {
		os.Setenv("no_proxy", noProxy)
		os.Setenv("NO_PROXY", noProxy)
		os.Setenv("http_proxy", httpProxy)
		os.Setenv("HTTP_PROXY", httpProxy)
		os.Setenv("HTTPS_PROXY", httpsProxy)
		os.Setenv("https_proxy", httpsProxy)
	}()

	testdata := map[string]string{
		"NO_PROXY":    "http://dummy.no.proxy",
		"http_proxy":  "http://dummy.http.proxy",
		"https_proxy": "http://dummy.https.proxy",
	}

	for k, v := range testdata {
		os.Setenv(k, v)
	}

	src := new(yaml.Container)
	dst := new(engine.Step)

	transformProxy()(dst, src, nil)

	for k, v := range testdata {
		if dst.Environment[k] != v {
			t.Errorf("Expect proxy varaible %s=%q, got %q", k, v, dst.Environment[k])
		}
	}
}
