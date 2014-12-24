package cfg

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/robfig/config"
	"github.com/voxelbrain/goptions"
)

const VERSION = "0.2"

type Options struct {
	Env     string        `goptions:"-e, --env, description='Application environment, defines config section'"`
	Rebuild bool          `goptions:"-r, --rebuild, description='Rebuild index only, do not start web server'"`
	RootDir string        `goptions:"-d, --root, description='Force set root dir'"`
	Help    goptions.Help `goptions:"-h, --help, description='Show this help'"`
}

var reader *config.Config
var options Options
var stdlogger *log.Logger
var instanceId string
var rootDir string

func String(key string) string {
	val, _ := reader.String(options.Env, key)
	return val
}

func Int(key string) int {
	val, _ := reader.Int(options.Env, key)
	return val
}

func Bool(key string) bool {
	val, _ := reader.Bool(options.Env, key)
	return val
}

func Log(msg string) {
	stdlogger.Printf("[LOG] %v | %s\n", time.Now().Format("2006/01/02 - 15:04:05"), msg)
}

func GetOptions() Options {
	return options
}

func GetInstanceId() string {
	return instanceId
}

func GetVersion() string {
	return VERSION
}

func GetRootDir() string {
	return rootDir
}

func init() {
	options = Options{
		Env:     "prod",
		Rebuild: false,
		RootDir: "",
	}
	goptions.ParseAndFail(&options)

	stdlogger = log.New(os.Stdout, "", 0)
	Log(fmt.Sprintf("Initializing application ver %s", VERSION))
	Log(fmt.Sprintf("Loading config with env set to %s", options.Env))

	if len(options.RootDir) > 0 {
		rootDir = options.RootDir
	} else {
		var err error
		if rootDir, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
			panic(err)
		}
		options.RootDir = rootDir
	}

	filePath := fmt.Sprintf("%s/%s.ini", GetRootDir(), options.Env)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		filePath = fmt.Sprintf("%s/config.ini", GetRootDir())
		Log(fmt.Sprintf("Loading config from %s", filePath))
		reader, _ = config.ReadDefault(filePath)
	} else {
		Log(fmt.Sprintf("Loading config from %s", filePath))
		reader, _ = config.ReadDefault(filePath)
	}

	hasher := md5.New()
	hasher.Write([]byte(time.Now().Format("2006/01/02 - 15:04:05")))
	instanceId = hex.EncodeToString(hasher.Sum(nil))[:5]
	Log(fmt.Sprintf("Initialized application instance ID %s", instanceId))
}
