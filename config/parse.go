package config

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Parse parses the configuration from bytes b.
func Parse(r io.Reader) (*Config, error) {
	out, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return ParseBytes(out)
}

// ParseBytes parses the configuration from bytes b.
func ParseBytes(b []byte) (*Config, error) {
	out := new(Config)
	err := yaml.Unmarshal(b, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ParseString parses the configuration from string s.
func ParseString(s string) (*Config, error) {
	return ParseBytes(
		[]byte(s),
	)
}

// ParseFile parses the configuration from path p.
func ParseFile(p string) (*Config, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

//
// parse multiple yaml documents delimited by ---
//

// ParseMulti parses the configurations from bytes b.
func ParseMulti(r io.Reader) ([]*Config, error) {
	var list []*Config
	scanner := bufio.NewScanner(r)
	row := 0
	buf := new(bytes.Buffer)
	for scanner.Scan() {
		row++
		txt := scanner.Text()
		if strings.HasPrefix(txt, "---") && row != 1 {
			out, err := Parse(buf)
			if err != nil {
				return nil, err
			}
			list = append(list, out)
			buf.Reset()
		} else {
			buf.WriteString(txt)
			buf.WriteByte('\n')
		}
	}
	out, err := Parse(buf)
	if err != nil {
		return nil, err
	}
	list = append(list, out)
	return list, nil
}

// ParseMultiBytes parses the configurations from bytes b.
func ParseMultiBytes(b []byte) ([]*Config, error) {
	return ParseMulti(
		bytes.NewBuffer(b),
	)
}

// ParseMultiString parses the configurations from string s.
func ParseMultiString(s string) ([]*Config, error) {
	return ParseMultiBytes(
		[]byte(s),
	)
}

// ParseMultiFile parses the configurations from path p.
func ParseMultiFile(p string) ([]*Config, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseMulti(f)
}
