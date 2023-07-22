package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/Luzifer/go_helpers/v2/backoff"
	"github.com/Luzifer/rconfig/v2"
)

var (
	cfg = struct {
		LogLevel         string        `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		MaxIterations    uint64        `flag:"max-iterations,i" vardefault:"maxIterations" description:"Maximum number of retries (0 = infinite)"`
		MaxIterationTime time.Duration `flag:"max-iteration-time" vardefault:"maxIterationTime" description:"How long to wait at most between iterations"`
		MaxTotalTime     time.Duration `flag:"max-total-time,t" vardefault:"maxTotalTime" description:"Deadline for overall executions (0 = infinite)"`
		MinIterationTime time.Duration `flag:"min-iteration-time" vardefault:"minIterationTime" description:"How long to wait before first retry"`
		Multiplier       float64       `flag:"mulitplier" vardefault:"mulitplier" description:"Mulitplier to apply to the wait-time after each retry (1.0 = constant backoff)"`
		Stdin            bool          `flag:"stdin" default:"false" description:"Pass stdin to command, to do so stdin will be fully buffered to memory before starting the command, enabling without input wil hang forever"`
		VersionAndExit   bool          `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	version = "dev"
)

func initApp() error {
	rconfig.SetVariableDefaults(map[string]string{
		"maxIterations":    strconv.FormatUint(backoff.DefaultMaxIterations, 10),
		"maxIterationTime": backoff.DefaultMaxIterationTime.String(),
		"maxTotalTime":     backoff.DefaultMaxTotalTime.String(),
		"minIterationTime": backoff.DefaultMinIterationTime.String(),
		"mulitplier":       strconv.FormatFloat(backoff.DefaultMultipler, 'f', -1, 64),
	})
	rconfig.AutoEnv(true)
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		return errors.Wrap(err, "parsing cli options")
	}

	l, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return errors.Wrap(err, "parsing log-level")
	}
	logrus.SetLevel(l)

	return nil
}

func main() {
	var err error
	if err = initApp(); err != nil {
		logrus.WithError(err).Fatal("initializing app")
	}

	if cfg.VersionAndExit {
		logrus.WithField("version", version).Info("backoff")
		os.Exit(0)
	}

	if len(rconfig.Args()[1:]) < 1 {
		logrus.Fatal("Usage: backoff [options] -- <command / args>")
	}

	var stdin io.ReadSeeker
	if cfg.Stdin {
		sbuf := new(bytes.Buffer)
		if _, err = io.Copy(sbuf, os.Stdin); err != nil {
			logrus.WithError(err).Fatal("reading stdin to buffer")
		}
		stdin = bytes.NewReader(sbuf.Bytes())
	}

	var (
		bo = backoff.NewBackoff().
			WithMaxIterations(cfg.MaxIterations).
			WithMaxIterationTime(cfg.MaxIterationTime).
			WithMaxTotalTime(cfg.MaxTotalTime).
			WithMinIterationTime(cfg.MinIterationTime).
			WithMultiplier(cfg.Multiplier)
		try int
	)

	if err = bo.Retry(func() error {
		try++

		logrus.WithField("try", try).Debug("starting execution")

		if cfg.Stdin {
			if _, err = stdin.Seek(0, io.SeekStart); err != nil {
				logrus.WithError(err).Fatal("resetting seek position in stdin buffer")
			}
		}

		cmd := exec.Command(rconfig.Args()[1], rconfig.Args()[2:]...)
		cmd.Env = os.Environ() // Env-Passthrough
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if cfg.Stdin {
			cmd.Stdin = stdin
		}

		return errors.Wrap(cmd.Run(), "executing command")
	}); err != nil {
		logrus.WithError(err).Fatal("retrying command")
	}
}
