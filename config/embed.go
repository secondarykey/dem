package config

type Embed struct {
	ID        string
	Kind      string
	Limit     int
	Namespace string
}

var currentEmbed *Embed

func SetCurrentEmbed(e *Embed) {
	currentEmbed = e
}

func GetCurrentEmbed() *Embed {
	return currentEmbed
}
