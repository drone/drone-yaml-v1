package linter

import (
	"testing"

	"github.com/drone/drone-yaml-v1/config"
)

func TestLint(t *testing.T) {
	testdata := `
pipeline:
  build:
    image: docker
    privileged: true
    network_mode: host
    volumes:
      - /tmp:/tmp
    commands:
      - go build
      - go test
  publish:
    image: plugins/docker
    repo: foo/bar
services:
  redis:
    image: redis
    entrypoint: [ /bin/redis-server ]
    command: [ -v ]
`

	conf, err := config.ParseString(testdata)
	if err != nil {
		t.Fatalf("Cannot unmarshal yaml %q. Error: %s", testdata, err)
	}
	if err := NewDefault(true).Lint(conf); err != nil {
		t.Errorf("Expected lint returns no errors, got %q", err)
	}
}

func TestLintErrors(t *testing.T) {
	testdata := []struct {
		from string
		want string
	}{
		{
			from: "",
			want: "Invalid or missing pipeline section",
		},
		//
		// custom volumes, networks
		//
		{
			from: "{ pipeline: { build: { image: 'golang' }  }, volumes: { custom: { driver: vieux/sshfs } } }",
			want: "Insufficient privileges to define custom volumes",
		},
		{
			from: "{ pipeline: { build: { image: 'golang' }  }, networks: { custom: { driver: overlay } } }",
			want: "Insufficient privileges to define custom networks",
		},
		//
		// pipeline containers
		//
		{
			from: "pipeline: { build: { image: '' }  }",
			want: "Invalid or missing image",
		},
		{
			from: "pipeline: { build: { image: golang, privileged: true }  }",
			want: "Insufficient privileges to use privileged mode",
		},
		{
			from: "pipeline: { build: { image: golang, shm_size: 10gb }  }",
			want: "Insufficient privileges to override shm_size",
		},
		{
			from: "pipeline: { build: { image: golang, dns: [ 8.8.8.8 ] }  }",
			want: "Insufficient privileges to use custom dns",
		},

		{
			from: "pipeline: { build: { image: golang, dns_search: [ example.com ] }  }",
			want: "Insufficient privileges to use dns_search",
		},
		{
			from: "pipeline: { build: { image: golang, devices: [ '/dev/tty0:/dev/tty0' ] }  }",
			want: "Insufficient privileges to use devices",
		},
		{
			from: "pipeline: { build: { image: golang, extra_hosts: [ 'somehost:162.242.195.82' ] }  }",
			want: "Insufficient privileges to use extra_hosts",
		},
		{
			from: "pipeline: { build: { image: golang, network_mode: host }  }",
			want: "Insufficient privileges to use network_mode",
		},
		{
			from: "pipeline: { build: { image: golang, networks: [ outside, default ] }  }",
			want: "Insufficient privileges to use networks",
		},
		{
			from: "pipeline: { build: { image: golang, volumes: [ '/opt/data:/var/lib/mysql' ] }  }",
			want: "Insufficient privileges to use volumes",
		},
		{
			from: "pipeline: { build: { image: golang, network_mode: 'container:name' }  }",
			want: "Insufficient privileges to use network_mode",
		},
		//
		// cannot override entypoint, command for script steps
		//
		{
			from: "pipeline: { build: { image: golang, commands: [ 'go build' ], entrypoint: [ '/bin/bash' ] } }",
			want: "Cannot configure both commands and entrypoint attributes",
		},
		{
			from: "pipeline: { build: { image: golang, commands: [ 'go build' ], command: [ '/bin/bash' ] } }",
			want: "Cannot override container command",
		},
		//
		// cannot override entypoint, command for plugin steps
		//
		{
			from: "pipeline: { publish: { image: plugins/docker, repo: foo/bar, entrypoint: [ '/bin/bash' ] } }",
			want: "Cannot override container entrypoint",
		},
		{
			from: "pipeline: { publish: { image: plugins/docker, repo: foo/bar, command: [ '/bin/bash' ] } }",
			want: "Cannot override container command",
		},
	}

	for _, test := range testdata {
		conf, err := config.ParseString(test.from)
		if err != nil {
			t.Fatalf("Cannot unmarshal yaml %q. Error: %s", test.from, err)
		}

		lerr := NewDefault(false).Lint(conf)
		if lerr == nil {
			t.Errorf("Expected lint error for configuration %q", test.from)
		} else if lerr.Error() != test.want {
			t.Errorf("Want error %q, got %q", test.want, lerr.Error())
		}
	}
}

var result error

func BenchmarkLint(b *testing.B) {
	var err error

	conf, err := config.ParseString(benchdata)
	if err != nil {
		b.Error(err)
		return
	}

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		err = NewDefault(false).Lint(conf)
		if err != nil {
			panic(err)
		}
	}
	result = err
}

var benchdata = `
workspace:
  base: /go
  path: src/github.com/drone/drone

pipeline:
  test:
    image: golang:1.8
    commands:
      - go get -u github.com/drone/drone-ui/dist
      - go get -u golang.org/x/tools/cmd/cover
      - go get -u golang.org/x/net/context
      - go get -u golang.org/x/net/context/ctxhttp
      - go get -u github.com/golang/protobuf/proto
      - go get -u github.com/golang/protobuf/protoc-gen-go
      - go test -cover $(go list ./... | grep -v /vendor/)

  test_postgres:
    image: golang:1.8
    environment:
      - DATABASE_DRIVER=postgres
      - DATABASE_CONFIG=host=postgres user=postgres dbname=postgres sslmode=disable
    commands:
      - go test github.com/drone/drone/store/datastore

  test_mysql:
    image: golang:1.8
    environment:
      - DATABASE_DRIVER=mysql
      - DATABASE_CONFIG=root@tcp(mysql:3306)/test?parseTime=true
    commands:
      - go test github.com/drone/drone/store/datastore

  build:
    image: golang:1.8
    commands: sh .drone.sh
    secrets: [ ssh_key ]
    when:
      event: [ push, tag ]

  publish_server_alpine:
    image: plugins/docker
    repo: drone/drone
    dockerfile: Dockerfile.alpine
    secrets: [ docker_username, docker_password ]
    tag: [ alpine ]
    when:
      branch: master
      event: push

  publish_server:
    image: plugins/docker
    repo: drone/drone
    secrets: [ docker_username, docker_password ]
    tag: [ latest ]
    when:
      branch: master
      event: push

  publish_agent_alpine:
    image: plugins/docker
    repo: drone/agent
    dockerfile: Dockerfile.agent.alpine
    secrets: [ docker_username, docker_password ]
    tag: [ alpine ]
    when:
      branch: master
      event: push

  publish_agent_arm:
    image: plugins/docker
    repo: drone/agent
    dockerfile: Dockerfile.agent.linux.arm
    secrets: [ docker_username, docker_password ]
    tag: [ linux-arm ]
    when:
      branch: master
      event: push

  publish_agent_arm64:
    image: plugins/docker
    repo: drone/agent
    dockerfile: Dockerfile.agent.linux.arm64
    secrets: [ docker_username, docker_password ]
    tag: [ linux-arm64 ]
    when:
      branch: master
      event: push

  publish_agent_amd64:
    image: plugins/docker
    repo: drone/agent
    dockerfile: Dockerfile.agent
    secrets: [ docker_username, docker_password ]
    tag: [ latest ]
    when:
      branch: master
      event: push

  release_server_alpine:
    image: plugins/docker
    repo: drone/drone
    dockerfile: Dockerfile.alpine
    secrets: [ docker_username, docker_password ]
    tag: [ 0.8-alpine ]
    when:
      event: tag

  release_agent_alpine:
    image: plugins/docker
    repo: drone/agent
    dockerfile: Dockerfile.agent.alpine
    secrets: [ docker_username, docker_password ]
    tag: [ 0.8-alpine ]
    when:
      event: tag

  release_server:
    image: plugins/docker
    repo: drone/drone
    secrets: [ docker_username, docker_password ]
    tag: [ 0.8, 0.8.1 ]
    when:
      event: tag

  release_agent:
    image: plugins/docker
    repo: drone/agent
    dockerfile: Dockerfile.agent
    secrets: [ docker_username, docker_password ]
    tag: [ 0.8, 0.8.1 ]
    when:
      event: tag

services:
  postgres:
    image: postgres:9.6
    environment:
      - POSTGRES_USER=postgres
  mysql:
    image: mysql:5.6.27
    environment:
      - MYSQL_DATABASE=test
      - MYSQL_ALLOW_EMPTY_PASSWORD=yes
`
