package confluence

//Confluence的备注宏，用于备注git的信息
const ConfluenceNoteMacro = `
<ac:structured-macro ac:name="note">
  <ac:parameter ac:name="icon">true</ac:parameter>
  <ac:parameter ac:name="title">Git提交信息</ac:parameter>
  <ac:rich-text-body>
    <p>%s</p>
  </ac:rich-text-body>
</ac:structured-macro>
`
