package resource

type ResourceProvider interface {
	GetResourceImage() string
	GetResourcePath() string
	GetResourceURLs() []string
}
