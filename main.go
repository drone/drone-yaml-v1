package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/drone/drone-yaml-v1/config"
	"github.com/drone/drone-yaml-v1/config/compiler"
	"github.com/drone/drone-yaml-v1/config/linter"
	"github.com/drone/drone-yaml-v1/version"

	"github.com/mattn/go-isatty"
)

var tty = isatty.IsTerminal(os.Stdout.Fd())

var (
	source       = kingpin.Arg("source", "source file location").Required().File()
	target       = kingpin.Arg("target", "target file location").String()
	trusted      = kingpin.Flag("trusted", "trusted mode").Bool()
	volume       = kingpin.Flag("volume", "attached volumes").Strings()
	network      = kingpin.Flag("network", "attached networks").Strings()
	environ      = kingpin.Flag("env", "environment variable").StringMap()
	images       = kingpin.Flag("privileged", "privileged images").Default("plugins/docker").Strings()
	base         = kingpin.Flag("base", "workspace base path").Default("/workspace").String()
	path         = kingpin.Flag("path", "wrokspace path").String()
	event        = kingpin.Flag("event", "event type").PlaceHolder("<event>").Enum("push", "pull_request", "tag", "deployment")
	repo         = kingpin.Flag("repo", "repository name").PlaceHolder("octocat/hello-world").String()
	branch       = kingpin.Flag("git-branch", "git commit branch").PlaceHolder("master").String()
	ref          = kingpin.Flag("git-ref", "git commit ref").PlaceHolder("refs/heads/master").String()
	deploy       = kingpin.Flag("deploy-to", "target deployment").PlaceHolder("production").String()
	platform     = kingpin.Flag("platform", "target platform").PlaceHolder("linux/amd64").String()
	secrets      = kingpin.Flag("secret", "secret variable").StringMap()
	registries   = kingpin.Flag("registry", "registry credentials").URLList()
	username     = kingpin.Flag("netrc-login", "netrc username").PlaceHolder("<token>").String()
	password     = kingpin.Flag("netrc-password", "netrc password").PlaceHolder("x-oauth-basic").String()
	machine      = kingpin.Flag("netrc-machine", "netrc machine").PlaceHolder("github.com").String()
	cpuset       = kingpin.Flag("cpu-set", "cpu set").PlaceHolder("0,1").String()
	cpushares    = kingpin.Flag("cpu-shares", "cpu shares").PlaceHolder("75").Int64()
	cpuquota     = kingpin.Flag("cpu-quota", "cpu quota").PlaceHolder("7500").Int64()
	memlimit     = kingpin.Flag("mem-limit", "memory limit").PlaceHolder("1GB").Bytes()
	memswaplimit = kingpin.Flag("mem-swap-limit", "memory swap limit").PlaceHolder("1GB").Bytes()
	shmsize      = kingpin.Flag("shmsize", "shmsize").PlaceHolder("1GB").Bytes()
)

func main() {
	kingpin.Version(version.Version.String())
	kingpin.Parse()

	conf, err := config.Parse(*source)
	if err != nil {
		log.Fatal(err)
	}

	if err := linter.NewDefault(*trusted).Lint(conf); err != nil {
		log.Fatal(err)
	}

	var secretList []compiler.Secret
	for k, v := range *secrets {
		secretList = append(secretList, compiler.Secret{
			Name:  k,
			Value: v,
		})
	}

	var registryList []compiler.Registry
	for _, uri := range *registries {
		if uri.User == nil {
			log.Fatalln("Expect registry format [user]:[password]@hostname")
		}
		password, ok := uri.User.Password()
		if !ok {
			log.Fatalln("Invalid or missing registry password")
		}
		registryList = append(registryList, compiler.Registry{
			Hostname: uri.Host,
			Username: uri.User.Username(),
			Password: password,
		})
	}

	var opts = []compiler.Option{
		compiler.WithClone(false),
		compiler.WithEnviron(*environ),
		compiler.WithLimits(
			compiler.Resources{
				CPUQuota:     *cpuquota,
				CPUShares:    *cpushares,
				CPUSet:       *cpuset,
				ShmSize:      int64(*shmsize),
				MemLimit:     int64(*memlimit),
				MemSwapLimit: int64(*memswaplimit),
			},
		),
		compiler.WithMetadata(
			compiler.Metadata{
				Branch:      *branch,
				Event:       *event,
				Ref:         *ref,
				Repo:        *repo,
				Platform:    *platform,
				Environment: *deploy,
			},
		),
		compiler.WithNetrc(*username, *password, *machine),
		compiler.WithNetworks(*network...),
		compiler.WithPrivileged(*images...),
		compiler.WithRegistry(registryList...),
		compiler.WithSecret(secretList...),
		compiler.WithVolumes(*volume...),
		compiler.WithWorkspace(*base, *path),
	}

	out, _ := compiler.New(opts...).Compile(conf)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	enc.Encode(out)
}
