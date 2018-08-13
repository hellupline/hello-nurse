package booru

type (
	NewClientFunc func(string, int) Client // nolint

	Client interface { // nolint: golint
		NewTag(string, int) Tag
		Name() string
	}

	Tag interface { // nolint: golint
		Fetch() error
		Posts() []Post

		Pages() int
	}

	Post interface { // nolint: golint
		Tags() []string
		Key() string
		Body() map[string]string
	}
)

var clientRegistry = map[string]NewClientFunc{}

func GetClient(name string) (NewClientFunc, bool) { // nolint: golint
	c, ok := clientRegistry[name]
	return c, ok
}

func RegisterClient(f NewClientFunc, name string) { // nolint: golint
	clientRegistry[name] = f
}
