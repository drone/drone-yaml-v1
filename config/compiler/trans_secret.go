package compiler

import (
	"strings"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/config"
	"github.com/drone/drone-yaml-v1/yaml"
)

func transformSecret(secrets ...Secret) Transform {
	lookup := map[string]Secret{}
	for _, secret := range secrets {
		lookup[strings.ToLower(secret.Name)] = secret
	}

	return func(dst *engine.Step, src *yaml.Container, _ *config.Config) {
		var injected []string
		for _, value := range src.Secrets.Secrets {
			secret, ok := lookup[strings.ToLower(value.Source)]
			if !ok {
				continue
			}

			if len(secret.Match) != 0 && !matchImage(dst.Image, secret.Match...) {
				continue
			}
			injected = append(injected, value.Target)
			dst.Secrets = append(dst.Secrets, &engine.Secret{
				Name:  strings.ToUpper(value.Target),
				Value: secret.Value,
				Mask:  true,
			})
		}
		if len(injected) == 0 {
			return
		}
		if dst.Environment == nil {
			dst.Environment = map[string]string{}
		}
		dst.Environment["DRONE_SECRETS"] = strings.Join(injected, ",")
	}
}
