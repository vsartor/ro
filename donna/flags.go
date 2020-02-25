package donna

var (
	globalFlags []string
	flags       []string
)

func init() {
	globalFlags = make([]string, 0)
	flags = make([]string, 0)
}

func HasGlobalFlag(name string) bool {
	for _, flag := range globalFlags {
		if flag == name {
			return true
		}
	}
	return false
}

func HasFlag(name string) bool {
	for _, flag := range flags {
		if flag == name {
			return true
		}
	}
	return false
}
