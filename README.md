# bucketctl

A simple CLI-tool written i Go for interacting with the BitBucket API.

Config blir lagret under `$HOME/.config/bucketctl/config.yaml`

## Installasjon

```shell
make install
```

## Help

`bucketctl` has a built-in help command

```shell 
bucketctl help
```

## Configuration

To use some of the functions you will need an access-token.
This token can be generated in your profile settings in Bitbucket.

For easy of use this token can be stored in a config file together with the rest of the configuration, e.g.

```shell
bucketctl config set --token <token> --base-url <url> --git-url <ssh>
```

It's also possible to create different contexts based on this base config.
To create a new context you can e.g. run

```shell
bucketctl config context create --context infra --key INFRA --include-repos
```

This lets you use the context `-x infra` instead of having to write `--project INFRA --include-repos`.

## Examples

### Basic

Fetch all projects

```shell
bucketctl get projects
```

Fetch all repositories under the `PROJ`-project

```shell
bucketctl get repos -p PROJ
```

To clone every repository in a project, run

```shell
bucketctl git clone -p PROJ
```

All repositories will then be cloned into a folder name `PROJ` unless another folder is specified.

To update the main brain of all repositories in a project, use the `--update flag`

```shell
bucketctl git clone -p PROJ --update
```

### Project Configuration

Run

```shell
bucketctl get project-config -p PROJ
```

to fetch the current project configuration.
New project configuration can be applied by running

```shell
bucketctl apply -f <PROJ>.yaml
```

To check which changes are to be made you can add the `--dry-run` flag.

Read the [project config documentation](./docs/PROJECT_CONFIG.md) for additional details.
