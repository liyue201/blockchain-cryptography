// Package book 承载所有随二进制打包的静态资源：
// 教材内容(content/)、示例代码(examples/)、前端构建产物(web/dist)。
// 三者均被 go:embed 打入单个二进制，部署时无需额外文件。
package book

import "embed"

// Content 教材文章（Markdown + 自定义 front matter）
//
//go:embed content
var Content embed.FS

// Examples 可运行的示例 Go 代码
//
//go:embed examples
var Examples embed.FS

// Web 前端构建产物（由 `npm run build` 生成到 web/dist）
//
//go:embed web/dist
var Web embed.FS
