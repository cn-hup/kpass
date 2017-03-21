package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/seccom/kpass/src"
	"github.com/seccom/kpass/src/appdir"
	"github.com/seccom/kpass/src/logger"
)

var (
	address  = flag.String("addr", "127.0.0.1:8088", `KPass service address to listen on`)
	bindHost = flag.String("bindhost", "", `Bind a host name to KPass client, default to service address`)
	dbPath   = flag.String("dbpath", "./kpass.kdb", `KPass database path - has to be a file, not directory`)
	logPath  = flag.String("logpath", "", `KPass log file path - has to be a file, not directory, default to stdout`)
	devMode  = flag.Bool("dev", false, "Development mode, will use memory database as default")
	certFile = flag.String("certFile", "", `certFile path, used to create TLS service, support HTTP/2`)
	keyFile  = flag.String("keyFile", "", `keyFile path, used to create TLS service, support HTTP/2`)
)

func main() {
	flag.Parse()
	if err := checkOSEnv(); err != nil {
		logger.Fatal(err)
	}
	if !strings.HasSuffix(*dbPath, ".kdb") {
		logger.Fatal(`Invalid dbpath, must have ".kdb" as extension.`)
	}

	if os.Getenv("APP_ENV") == "" {
		if *devMode {
			os.Setenv("APP_ENV", "development")
		} else {
			os.Setenv("APP_ENV", "production")
		}
	}

	if *bindHost == "" {
		*bindHost = *address
	}
	if !strings.HasPrefix(*bindHost, "http") {
		*bindHost = "http://" + *bindHost
	}
	if *certFile != "" && *keyFile != "" && !strings.HasPrefix(*bindHost, "https") {
		*bindHost = "https" + (*bindHost)[4:]
	}

	var state int32 = 1
	app, db := src.New(*dbPath, *bindHost)
	ac := make(chan error)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		if *certFile != "" && *keyFile != "" {
			ac <- app.ListenTLS(*address, *certFile, *keyFile)
		} else {
			ac <- app.Listen(*address)
		}
	}()

	go func() {
		time.Sleep(600 * time.Millisecond)
		logger.Info("Start KPass: " + *bindHost)
		startBrowser(*bindHost)
	}()

	select {
	case err := <-ac:
		if err != nil && atomic.LoadInt32(&state) == 1 {
			logger.Err(err)
		}
	case sig := <-sc:
		atomic.StoreInt32(&state, 0)
		logger.Info(fmt.Sprintf("Got signal [%d] to exit.", sig))
		if err := app.Close(); err != nil {
			logger.Err(err)
		}
	}

	if err := db.Close(); err != nil {
		logger.Err(err)
	}
	os.Exit(int(atomic.LoadInt32(&state)))
}

// startBrowser tries to open the URL in a browser
// and reports whether it succeeds.
func startBrowser(url string) bool {
	// try to start the browser
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}

func checkOSEnv() error {
	app, err := os.Executable()
	if err != nil {
		return err
	}
	if strings.HasSuffix(app, "KPass.app/Contents/MacOS/kpass") { // run as mac os app
		path := appdir.General("KPass")
		if err := os.MkdirAll(path, os.ModeDir|os.FileMode(0755)); err != nil {
			return err
		}
		*dbPath = filepath.Join(path, "kpass.kdb")

		path = appdir.Logs("KPass")
		if err := os.MkdirAll(path, os.ModeDir|os.FileMode(0755)); err != nil {
			return err
		}
		*logPath = filepath.Join(path, "kpass.log")
	}
	return nil
}
