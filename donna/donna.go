package donna

// Donna handles command line arguments for Ro.

// It handles two different types of arguments: options and arguments.

// Options are optional command line arguments with default values and
// that starts with "-", either containing a parameter if integer or
// string, or simply being a flag.

// Arguments are command line arguments that do not have a preceding dash.

// There are two kinds of options, global options, handled at startup and
// consumed by ro's entry point, and command specific options which are
// handled by the executing code.
