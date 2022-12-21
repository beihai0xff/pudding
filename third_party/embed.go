package third_party

import (
	"embed"
)

//go:embed OpenAPI/*
var Embed embed.FS
