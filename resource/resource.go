package resource

type ResourceProvider interface {
	GetResourcePath() string
	GetResourceURLs() []string
}
