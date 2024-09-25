# SSM-Shell

> connect to ECS containers or EC2 Instances with ease.

![Example of SSM-Shell](/examples/example.png)

## Usage

Download your binary from [here](https://github.com/digitalkaoz/ssm-shell/releases) 

simply prepend the binary call with your AWS environment, e.g.

```shell
$ AWS_PROFILE=xyz AWS_REGION=eu-central-1 ssm-shell

```
## Development

```shell
$ go get
$ AWS_PROFILE=xyz AWS_REGION=eu-central-1 go run .
```

To build the binary locally:

```shell
$ goreleaser build --snapshot --clean
$ cp ./dist/ssm_shell_xxx/ssm-shell /usr/local/bin/ssm-shell
$ chmod a+x /usr/local/bin/ssm-shell
```

## TODO

- [ ] get rid of shellout to aws cli, session-manager might be used directed
- [ ] tests
- [ ] better view templates
- [ ] add fuzzy cli args for selecting things
- [ ] region switch?
