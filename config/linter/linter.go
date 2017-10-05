package linter

import "github.com/drone/drone-yaml-v1/config"

// New returns a new Linter that executes checks sequentially.
func New(checks ...Check) *Linter {
	return &Linter{checks}
}

// NewDefault returns a new Linter that executes default checks.
func NewDefault(trusted bool) *Linter {
	return New(
		CheckPipeline,
		CheckContainer(CheckCommand),
		CheckContainer(CheckCommands),
		CheckContainer(CheckEntrypoint),
		CheckContainer(CheckImage),
		CheckTrusted(trusted),
		CheckVolumes(trusted),
		CheckNetworks(trusted),
	)
}

// A Linter lints a pipeline configuration.
type Linter struct {
	checks []Check
}

// Lint evaluates the linter rules against the given configuration.
func (l *Linter) Lint(conf *config.Config) error {
	for _, check := range l.checks {
		if err := check(conf); err != nil {
			return err
		}
	}
	return nil
}
