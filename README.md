# Confluence SDK
[![GoDoc](https://godoc.org/github.com/go-http/confluence?status.svg)](https://godoc.org/github.com/go-http/confluence)
[![Build Status](https://travis-ci.com/go-http/confluence.svg?branch=master)](https://travis-ci.com/go-http/confluence)

## Usage

```golang
import (
	"github.com/go-http/confluence"
)
```

具体用法可以查看[cli目录](cli/)下的两个范例：

- [cli/sync2confluence](cli/sync2confluence/): 用于将指定目录同步到Confluence空间。
- [cli/confluence_exporter](cli/confluence_exporter/): 用于将Confluence空间导出到指定目录。
