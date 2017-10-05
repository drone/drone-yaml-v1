The drone-yaml-v1 package provides a parser, linter and compiler for version 1 of the drone configuration file. The compiler converts the yaml configuration file to an intermediate representation suitable for [drone-runtime](https://github.com/drone/drone-runtime) execution.

Compile the yaml configuration file to the intermediate representation:

```text
drone-yaml samples/1_simple.yml samples/1_simple.json
```

Execute the intermediate representation using the runtime tools:

```text
drone-runtime samples/1_simple.json
```
