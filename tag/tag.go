package tag

type TagDef interface {
	ToFootTag() string
	ToHeadTag() string
	ToTag() string
}
