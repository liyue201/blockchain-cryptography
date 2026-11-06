// Package book carries all static resources bundled with the binary:
// teaching content (content/), example code (examples/), and the frontend build output (web/dist).
// All three are embedded into a single binary via go:embed, so no extra files are needed at deploy time.
package book

import "embed"

// Content teaching articles (Markdown + custom front matter)
//
//go:embed content
var Content embed.FS

// Examples runnable example Go code
//
//go:embed examples
var Examples embed.FS

// Web frontend build output (generated to web/dist by `npm run build`)
//
//go:embed web/dist
var Web embed.FS
