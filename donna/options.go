package donna

var (
	globalOptions map[string]string
	options map[string]string
)

func init() {
	globalOptions = make(map[string]string)
	options = make(map[string]string)
}

func GetGlobalOption(name string) (string, bool) {
	value, exists := globalOptions[name]
	return value, exists
}

func GetOption(name string) (string, bool) {
	value, exists := options[name]
	return value, exists
}
