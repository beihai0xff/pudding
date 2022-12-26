package third_party

import (
	"embed"
)

//go:embed swagger-ui/*
var Embed embed.FS
