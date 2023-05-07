package data

type Provider interface {
	Categories() []Category
}

type Category interface {
	Title() string
	Links() []Link
}

type Link interface {
	Title() string
	URL() string
	Description() string
}
