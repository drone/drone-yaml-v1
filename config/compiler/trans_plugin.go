package compiler

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml-v1/config"
	"github.com/drone/drone-yaml-v1/yaml"
)

// transformPlugin is an internal transform used to convert plugin
// parameters to environment variables.
func transformPlugin(dst *engine.Step, src *yaml.Container, _ *config.Config) {
	if dst.Environment == nil {
		dst.Environment = map[string]string{}
	}
	paramsToEnv(src.Vargs, dst.Environment)
}

// paramsToEnv uses reflection to convert a map[string]interface to a list
// of environment variables.
func paramsToEnv(from map[string]interface{}, to map[string]string) error {
	for k, v := range from {
		if v == nil {
			continue
		}

		t := reflect.TypeOf(v)
		vv := reflect.ValueOf(v)

		k = "PLUGIN_" + strings.ToUpper(k)

		switch t.Kind() {
		case reflect.Bool:
			to[k] = strconv.FormatBool(vv.Bool())

		case reflect.String:
			to[k] = vv.String()

		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
			to[k] = fmt.Sprintf("%v", vv.Int())

		case reflect.Float32, reflect.Float64:
			to[k] = fmt.Sprintf("%v", vv.Float())

		case reflect.Map:
			to[k] = encodeMap(vv.Interface())

		case reflect.Slice:
			to[k] = encodeSlice(vv.Interface())
		}
	}
	return nil
}
