# Luzifer / backoff

`backoff` is a small CLI util wrapping [`github.com/Luzifer/go_helpers/v2/backoff`](https://pkg.go.dev/github.com/Luzifer/go_helpers/v2@v2.20.0/backoff) to be used in shell scripts

## Usage

```console
# backoff --help
Usage of backoff:
      --log-level string              Log level (debug, info, warn, error, fatal) (default "info")
      --max-iteration-time duration   How long to wait at most between iterations (default 1m0s)
  -i, --max-iterations uint           Maximum number of retries (0 = infinite)
  -t, --max-total-time duration       Deadline for overall executions (0 = infinite)
      --min-iteration-time duration   How long to wait before first retry (default 100ms)
      --mulitplier float              Mulitplier to apply to the wait-time after each retry (1.0 = constant backoff) (default 1.5)
      --stdin                         Pass stdin to command, to do so stdin will be fully buffered to memory before starting the command, enabling without input wil hang forever
      --version                       Prints current version and exits

# backoff -i 10 --log-level=debug -- false
time="2023-07-22T14:38:28+02:00" level=debug msg="starting execution" try=1
time="2023-07-22T14:38:28+02:00" level=debug msg="starting execution" try=2
time="2023-07-22T14:38:28+02:00" level=debug msg="starting execution" try=3
time="2023-07-22T14:38:28+02:00" level=debug msg="starting execution" try=4
time="2023-07-22T14:38:28+02:00" level=debug msg="starting execution" try=5
time="2023-07-22T14:38:29+02:00" level=debug msg="starting execution" try=6
time="2023-07-22T14:38:30+02:00" level=debug msg="starting execution" try=7
time="2023-07-22T14:38:31+02:00" level=debug msg="starting execution" try=8
time="2023-07-22T14:38:33+02:00" level=debug msg="starting execution" try=9
time="2023-07-22T14:38:35+02:00" level=debug msg="starting execution" try=10
time="2023-07-22T14:38:35+02:00" level=fatal msg="retrying command" error="Maximum iterations reached: executing command: exit status 1"
```
