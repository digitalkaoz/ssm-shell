# SSM-Shell

> connect to ECS containers or EC2 Instances with ease.

## Usage

Download binary from here:

simply prepend the binary call with your AWS environment, e.g.

```shell
$ AWS_PROFILE=xyz AWS_REGION=eu-central-1 ecs-shell

```
## Development

```shell
$ go get
$ AWS_PROFILE=xyz AWS_REGION=eu-central-1 go run .
```

To build and distribute the binary:

$ goreleaser build --snapshot --clean
$ cp ./dist/ssm_shell_xxx/ssm-shell /usr/local/bin/ssm_shell
$ chmod a+x /usr/local/bin/ssm_shell


## TODO

- [ ] get rid of shellout to aws cli, session-manager might be used directed
- [ ] tests
- [ ] better view templates
- [ ] add fuzzy cli args for selecting things
- [ ] region switch?