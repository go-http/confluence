package confluence

//Confluence的备注宏，用于备注git的信息
const ConfluenceNoteMacro = `
<br />
<br />
<br />
<hr />
<ac:structured-macro ac:name="note">
  <ac:parameter ac:name="icon">true</ac:parameter>
  <ac:parameter ac:name="title">修改历史</ac:parameter>
  <ac:rich-text-body>
    <p>%s</p>
  </ac:rich-text-body>
</ac:structured-macro>
`

const NewConfluenceNoteMacro = `
<br />
<br />
<br />
<hr />
<ac:structured-macro ac:name="note">
  <ac:parameter ac:name="icon">true</ac:parameter>
  <ac:parameter ac:name="title">修改历史</ac:parameter>
  <ac:rich-text-body>
    {{range .CommitList}}
		<p>{{.}}</p>
	{{end}}
  <br />
  <p>渲染自<a href="{{.GitUrl}}">{{.GitUrl}}</a>仓库的<a href="{{.FileUrl}}}">{{.FileName}}</a>文件</p>
  </ac:rich-text-body>
</ac:structured-macro>
`

var ConfluenceNoteSplite = `
<br />
<br />
<br />
<hr />
`
