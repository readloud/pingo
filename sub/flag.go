package sub

import (
	"flag"
	"fmt"
	"os"
	"regexp"
)

type Flag struct {
	Count        int
	Data         string
	Exploit      string
	Help         bool
	Interval     string
	NoColor      bool
	Target       string
	TTL          int
	Unprivileged bool
	UseIPv4      bool
	UseIPv6      bool
	Verbose      bool
	Version      bool
}

var version = "v0.1.3"

var usageExploit = `pingo attack

*It is for educational purposes or pentesting against your own server.
*Don't use it to attack someone else's server.

USAGE:
  -x flood	Ping flood
  -x land	LAND attack
  -x pod	Ping of Death
  -x smurf	Smurf attack

EXAMPLES:
  pingo -x pod example.com`

// Parse flags
func (f *Flag) Parse() error {
	flag.BoolVar(&f.NoColor, "no-color", false, "disable coloring of outputs")
	flag.IntVar(&f.Count, "c", 0, "ping <count> times")
	flag.StringVar(&f.Data, "d", "PINGO", "custom data string")
	flag.BoolVar(&f.Help, "h", false, "print usage")
	flag.StringVar(&f.Interval, "i", "1", "interval per ping")
	flag.IntVar(&f.TTL, "t", 64, "set TTL (time to live) of the packet")
	flag.BoolVar(&f.Unprivileged, "u", false, "unprivileged (UDP) ping")
	flag.BoolVar(&f.Verbose, "v", false, "verbose mode")
	flag.StringVar(&f.Exploit, "x", "", "exploit with some attack *under development so cannot use it yet")
	flag.BoolVar(&f.UseIPv4, "4", true, "use IPv4")
	flag.BoolVar(&f.UseIPv6, "6", false, "use IPv6")
	flag.Parse()

	if f.Help || (len(os.Args) == 2 && os.Args[1] == "help") {
		flag.Usage()
		os.Exit(0)
	} else if f.Version || (len(os.Args) == 2 && os.Args[1] == "version") {
		fmt.Printf("pingo %s\n", version)
		os.Exit(0)
	}

	if len(os.Args) == 1 {
		flag.Usage()
		return ErrNoArguments
	}

	// validate arguments values
	if containsSlice(flag.Args(), "x") && f.Exploit == "" {
		return ErrNoArgumentValue
	}

	// validate an exploit
	if f.Exploit != "" && !validExploit(f.Exploit) {
		fmt.Println(usageExploit)
		return ErrInvalidAttackType
	}

	// validate interval
	if !validInterval(f.Interval) {
		return ErrIncorrectValueInterval
	}

	f.Target = flag.Arg(0)

	return nil
}

// Validate exploit
func validExploit(exploit string) bool {
	switch exploit {
	case "flood", "land", "pod", "smurf":
		return true
	default:
		return false
	}
}

// Validate interval
func validInterval(interval string) bool {
	r, _ := regexp.Compile(`^([1-9][0-9]*|0)`)
	return r.MatchString(interval)
}
