package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/kr/fs"
	glob "github.com/obeattie/ohmyglob"
	yaml "gopkg.in/yaml.v2"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: goglobls [options] [path...]\n")
	flag.PrintDefaults()
}

// Config Configuration structure
type Config struct {
	Types []struct {
		Name  string   `yaml:"name"`
		Paths []string `yaml:"paths"`
	}
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	flag.Usage = usage
	configFile := flag.String("config", "", "Config file (required)")
	ignoreCase := flag.Bool("i", false, "Ignore case (Unix Only)")
	verbose := flag.Bool("v", false, "Verbose")

	var types arrayFlags
	flag.Var(&types, "type", "Type to filter (add more than one by repeating the option)")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		usage()
		fmt.Fprintln(os.Stderr, "Path is missing.")
		os.Exit(1)
	}

	if *configFile == "" {
		usage()
		fmt.Fprintln(os.Stderr, "Config file is missing.")
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(*configFile)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	config := Config{}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	patterns := make([]string, 0)

	for _, _type := range config.Types {
		if len(types) > 0 {
			foundType, _ := inArray(_type.Name, types)
			if !foundType {
				continue
			}
		}
		for _, pattern := range _type.Paths {
			if runtime.GOOS == "windows" {
				patterns = append(patterns, strings.ToLower(pattern))
			} else {
				if *ignoreCase == true {
					patterns = append(patterns, strings.ToLower(pattern))
				} else {
					patterns = append(patterns, pattern)
				}
			}
		}
	}

	if *verbose {
		fmt.Printf("Patterns:\n")
		fmt.Printf("%v\n", patterns)
	}

	set, _ := glob.CompileGlobSet(patterns, glob.DefaultOptions)

	walker := fs.Walk(args[0])

	for walker.Step() {
		if err := walker.Err(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		var target string

		if runtime.GOOS == "windows" {
			target = strings.ToLower(strings.Replace(walker.Path(), "\\", "/", -1))
		} else {
			if *ignoreCase == true {
				target = strings.ToLower(walker.Path())
			} else {
				target = walker.Path()
			}
		}

		match := walker.Stat().Mode().IsRegular() && set.MatchString(target)

		if match && !*verbose {
			fmt.Printf("%v\n", walker.Path())
		} else if *verbose {
			fmt.Printf("Match of %v is %v\n", target, match)
		}
	}
}

func inArray(val string, array []string) (exists bool, index int) {
	exists = false
	index = -1

	for i, v := range array {
		if val == v {
			index = i
			exists = true
			return
		}
	}

	return
}
