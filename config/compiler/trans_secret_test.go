package compiler

import (
	"testing"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/yaml"
)

func Test_transformSecret(t *testing.T) {
	testdatum := []struct {
		image   string
		source  string
		target  string
		name    string
		value   string
		match   []string
		matched bool
	}{
		{
			image:   "golang",
			source:  "docker_password",
			target:  "DOCKER_PASSWORD",
			name:    "docker_password",
			value:   "pa55word",
			matched: true,
		},
		{
			image:   "golang",
			source:  "username",
			target:  "DOCKER_USERNAME",
			name:    "username",
			value:   "octocat",
			matched: true,
		},
		// case insensitive
		{
			image:   "golang",
			source:  "docker_email",
			target:  "DOCKER_EMAIL",
			name:    "DOCKER_EMAIL",
			value:   "octocat@github.com",
			matched: true,
		},
		// the secret source does not match the name
		{
			image:   "golang",
			source:  "aws_secret_access_key",
			target:  "aws_secret_access_key",
			name:    "aws_access_key_id",
			value:   "1234",
			matched: false,
		},
		// the secret source matches the name, but does
		// not match the image whitelist
		{
			image:   "golang",
			source:  "heroku_token",
			target:  "heroku_token",
			name:    "heroku_token",
			value:   "1234",
			match:   []string{"node"},
			matched: false,
		},
		// the secret source matches the name AND matches
		// the image in the whitelist.
		{
			image:   "golang",
			source:  "github_token",
			target:  "GITHUB_TOKEN",
			name:    "github_token",
			value:   "1234",
			match:   []string{"node", "golang"},
			matched: true,
		},
	}

	for _, testdata := range testdatum {
		src := new(yaml.Container)
		src.Secrets.Secrets = []*yaml.Secret{
			{
				Source: testdata.source,
				Target: testdata.target,
			},
		}

		dst := new(engine.Step)
		dst.Image = testdata.image

		secret := Secret{
			Name:  testdata.name,
			Value: testdata.value,
			Match: testdata.match,
		}

		transformSecret(secret)(dst, src, nil)

		if !testdata.matched {
			if len(dst.Secrets) != 0 {
				t.Errorf("Expect container not granted access to secret [%s]",
					testdata.name,
				)
			}
			continue
		}

		if testdata.matched && len(dst.Secrets) == 0 {
			t.Errorf("Expect container granted access to secret [%s]",
				testdata.name,
			)
			continue
		}

		if got, want := dst.Secrets[0].Value, testdata.value; got != want {
			t.Errorf("Expect secret value %s, got %s", want, got)
		}
		if got, want := dst.Secrets[0].Name, testdata.target; got != want {
			t.Errorf("Expect secret name %s, got %s", want, got)
		}
		if got, want := dst.Environment["DRONE_SECRETS"], testdata.target; got != want {
			t.Errorf("Expect DRONE_SECRETS=%s, got %s", want, got)
		}
	}
}
