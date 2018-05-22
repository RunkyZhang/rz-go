package global

import "flag"

var (
	Arguments = getArguments()
)

func getArguments() (map[string]string) {
	flag.String(ArgumentNameConfig, "application.json", "config file")
	flag.Parse()

	arguments := map[string]string{}
	arguments[ArgumentNameConfig] = flag.Lookup(ArgumentNameConfig).Value.String()

	return arguments
}
