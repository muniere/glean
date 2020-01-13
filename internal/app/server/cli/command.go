package cli

import (
	"io/ioutil"
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	pubsub "github.com/muniere/glean/internal/app/server/pubsub/manager"
	"github.com/muniere/glean/internal/pkg/lumber"
	"github.com/muniere/glean/internal/pkg/pathname"
	"github.com/muniere/glean/internal/pkg/rpc"
	"github.com/muniere/glean/internal/pkg/signals"
	"github.com/muniere/glean/internal/pkg/std"
	"github.com/muniere/glean/internal/pkg/task"
)

const (
	cmdLogName = "glean.cmd.log"
	outLogName = "glean.out.log"
	errLogName = "glean.err.log"
)

func NewCommand() *cobra.Command {
	return assemble(&cobra.Command{
		Use: "gleand",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, args)
		},
	})
}

type context struct {
	options optionSet
}

type optionSet struct {
	address     string
	port        int
	parallel    int
	concurrency int
	minWidth    int
	maxWidth    int
	minHeight   int
	maxHeight   int
	overwrite   bool
	dataDir     string
	logDir      string
	dryRun      bool
	verbose     bool
}

func assemble(cmd *cobra.Command) *cobra.Command {
	flags := cmd.Flags()
	flags.String("address", rpc.LocalAddr, "Address to bind")
	flags.Int("port", rpc.Port, "Port to bind")
	flags.Int("parallel", task.Parallel, "The number of workers for download")
	flags.Int("concurrency", task.Concurrency, "Concurrency of download tasks per worker")
	flags.Int("min-width", -1, "Minimum width of images")
	flags.Int("max-width", -1, "Maximum width of images")
	flags.Int("min-height", -1, "Minimum height of images")
	flags.Int("max-height", -1, "Maximum height of images")
	flags.String("data-dir", "", "Base directory to download files")
	flags.String("log-dir", "", "Path to log directory")
	flags.BoolP("dry-run", "n", false, "Do not perform actions actually")
	flags.BoolP("verbose", "v", false, "Show verbose messages")
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	ctx, err := parse(args, cmd.Flags())
	if err != nil {
		return err
	}

	if err := prepare(ctx); err != nil {
		return err
	}

	lumber.Info(std.NewDict(std.Pair("module", "root"), std.Pair("event", "start"), std.Pair("pid", os.Getpid())))

	// build
	manager := pubsub.NewManager(translate(ctx.options))

	// start
	if err = manager.Start(); err != nil {
		lumber.Fatal(std.NewDict(std.Pair("module", "root"), std.Pair("event", "start::error"), std.Pair("error", err)))
	}

	// wait
	sigs := []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	lumber.Info(std.NewDict(std.Pair("module", "root"), std.Pair("event", "signal::wait"), std.Pair("signals", signals.Join(sigs, ", "))))

	sig := signals.Wait(sigs...)
	lumber.Info(std.NewDict(std.Pair("module", "root"), std.Pair("event", "signal::recv"), std.Pair("signal", sig.String())))

	// stop
	if err := manager.Stop(); err != nil {
		lumber.Fatal(std.NewDict(std.Pair("module", "root"), std.Pair("event", "stop::error"), std.Pair("error", err)))
	}

	lumber.Info(std.NewDict(std.Pair("module", "root"), std.Pair("event", "stop"), std.Pair("pid", os.Getpid())))

	return nil
}

func parse(args []string, flags *pflag.FlagSet) (context, error) {
	optionSet, err := decode(flags)
	if err != nil {
		return context{}, err
	}

	ctx := context{
		options: optionSet,
	}
	return ctx, nil
}

func decode(flags *pflag.FlagSet) (optionSet, error) {
	address, err := flags.GetString("address")
	if err != nil {
		return optionSet{}, err
	}

	port, err := flags.GetInt("port")
	if err != nil {
		return optionSet{}, err
	}

	parallel, err := flags.GetInt("parallel")
	if err != nil {
		return optionSet{}, err
	}

	concurrency, err := flags.GetInt("concurrency")
	if err != nil {
		return optionSet{}, err
	}

	minWidth, err := flags.GetInt("min-width")
	if err != nil {
		return optionSet{}, err
	}

	maxWidth, err := flags.GetInt("max-width")
	if err != nil {
		return optionSet{}, err
	}

	minHeight, err := flags.GetInt("min-height")
	if err != nil {
		return optionSet{}, err
	}

	maxHeight, err := flags.GetInt("max-height")
	if err != nil {
		return optionSet{}, err
	}

	dataDir, err := flags.GetString("data-dir")
	if err != nil {
		return optionSet{}, err
	}

	logDir, err := flags.GetString("log-dir")
	if err != nil {
		return optionSet{}, err
	}

	dryRun, err := flags.GetBool("dry-run")
	if err != nil {
		return optionSet{}, err
	}

	verbose, err := flags.GetBool("verbose")
	if err != nil {
		return optionSet{}, err
	}

	opts := optionSet{
		address:     address,
		port:        port,
		parallel:    parallel,
		concurrency: concurrency,
		minWidth:    minWidth,
		maxWidth:    maxWidth,
		minHeight:   minHeight,
		maxHeight:   maxHeight,
		dataDir:     dataDir,
		logDir:      logDir,
		dryRun:      dryRun,
		verbose:     verbose,
	}

	return opts, nil
}

func prepare(ctx context) error {
	if ctx.options.verbose {
		logrus.SetLevel(logrus.TraceLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	if len(ctx.options.logDir) == 0 {
		if err := prepareForConsoleLog(ctx); err != nil {
			return err
		}
		base := &logrus.TextFormatter{
			DisableColors:    false,
			DisableTimestamp: false,
			FullTimestamp:    true,
			TimestampFormat:  "15:04:05.000",
		}
		logrus.SetFormatter(&lumber.TextFormatter{base})
	} else {
		if err := prepareForFileLog(ctx); err != nil {
			return err
		}
		base := &logrus.JSONFormatter{
			TimestampFormat:  "15:04:05.000",
			DisableTimestamp: false,
			DataKey:          "fields",
			FieldMap:         nil,
			CallerPrettyfier: nil,
			PrettyPrint:      false,
		}
		logrus.SetFormatter(&lumber.JSONFormatter{base})
	}

	return nil
}

func prepareForConsoleLog(ctx context) error {
	logrus.SetOutput(os.Stderr)
	return nil
}

func prepareForFileLog(ctx context) error {
	var err error

	logrus.SetOutput(ioutil.Discard)

	err = prepareCmdFileLog(ctx)
	if err != nil {
		return err
	}

	err = prepareOutFileLog(ctx)
	if err != nil {
		return err
	}

	err = prepareErrFileLog(ctx)
	if err != nil {
		return err
	}

	return nil
}

func prepareCmdFileLog(ctx context) error {
	file, err := os.OpenFile(
		pathname.Join(ctx.options.logDir, cmdLogName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644,
	)
	if err != nil {
		return err
	}

	logrus.AddHook(lumber.NewFileHookWithFilter(file, logrus.AllLevels, func(entry *logrus.Entry) bool {
		return entry.Data["command"] != nil
	}))
	return nil
}

func prepareOutFileLog(ctx context) error {
	file, err := os.OpenFile(
		pathname.Join(ctx.options.logDir, outLogName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644,
	)
	if err != nil {
		return err
	}

	logrus.AddHook(lumber.NewFileHook(file, logrus.AllLevels))
	return nil
}

func prepareErrFileLog(ctx context) error {
	file, err := os.OpenFile(
		pathname.Join(ctx.options.logDir, errLogName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644,
	)
	if err != nil {
		return err
	}

	logrus.AddHook(lumber.NewFileHook(file, []logrus.Level{
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
	}))
	return nil
}

func translate(options optionSet) pubsub.Config {
	return pubsub.Config{
		Address:     options.address,
		Port:        options.port,
		DataDir:     options.dataDir,
		Parallel:    options.parallel,
		Concurrency: options.concurrency,
		MinWidth:    options.minWidth,
		MaxWidth:    options.maxWidth,
		MinHeight:   options.minHeight,
		MaxHeight:   options.maxHeight,
		Overwrite:   options.overwrite,
		LogDir:      options.logDir,
		DryRun:      options.dryRun,
		Verbose:     options.verbose,
	}
}
