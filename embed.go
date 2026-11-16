// Package padtemplate exposes the canonical starter files to the CLI.
package padtemplate

import "embed"

// TemplateFS contains the versioned project starter structure.
//
//go:embed template/* template/.gitignore template/.pad-template.yaml template/.agent/**
var TemplateFS embed.FS
