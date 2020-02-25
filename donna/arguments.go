package donna

type Args struct {
	numConsumed int
	arguments   []string
}

var args Args

func init() {
	args.arguments = make([]string, 0)
}

// Returns the next command line argument in the chain.
func NextArg() (string, bool) {
	// Make sure we've not consumed all arguments
	if args.numConsumed == len(args.arguments) {
		return "", false
	}

	arg := args.arguments[args.numConsumed]
	args.numConsumed++
	return arg, true
}
