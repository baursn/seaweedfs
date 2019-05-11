package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/chrislusf/seaweedfs/weed/command"
	"github.com/chrislusf/seaweedfs/weed/glog"
	"github.com/kardianos/osext"
	"github.com/jacobsa/daemonize"
)

var (
	options      = flag.String("o", "", "comma separated options rw,uid=xxx,gid=xxx")
	isForeground = flag.Bool("foreground", false, "starts as a daemon")
)

func main() {

	flag.Parse()

	device := flag.Arg(0)
	mountPoint := flag.Arg(1)

	fmt.Printf("source: %v\n", device)
	fmt.Printf("target: %v\n", mountPoint)

	maybeSetupPath()

	if !*isForeground {
		startAsDaemon()
		return
	}

	parts := strings.SplitN(device, "/", 2)
	filer, filerPath := parts[0], parts[1]

	command.RunMount(
		filer, "/"+filerPath, mountPoint, "", "000", "",
		4, true, 0, 1000000)

}

func maybeSetupPath() {
	// sudo mount -av may not include PATH in some linux, e.g., Ubuntu
	hasPathEnv := false
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "PATH=") {
			hasPathEnv = true
		}
		fmt.Println(e)
	}
	if !hasPathEnv {
		os.Setenv("PATH", "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")
	}
}

func startAsDaemon() {

	// adapted from gcsfuse

	// Find the executable.
	var path string
	path, err := osext.Executable()
	if err != nil {
		glog.Fatalf("osext.Executable: %v", err)
	}

	// Set up arguments. Be sure to use foreground mode.
	args := append([]string{"-foreground"}, os.Args[1:]...)

	// Pass along PATH so that the daemon can find fusermount on Linux.
	env := []string{
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
	}

	err = daemonize.Run(path, args, env, os.Stdout)
	if err != nil {
		glog.Fatalf("daemonize.Run: %v", err)
	}

}
