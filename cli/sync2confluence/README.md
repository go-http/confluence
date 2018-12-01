
# Sync2Confluence

`Sync2Confluence`是用于同步本地目录到Confluence空间的命令行工具。

## 安装
1. 安装`go get github.com/athurg/go-confluence/cli/sync2confluence`
1. 执行`sync2confluence -addr WIKI地址 -s 空间标识 -u 用户名 -p 密码 -d 文档路径`
1. 完整的参数信息可以执行`sync2confluence -h`查看

## 目录树和Confluence空间的关系

Confluence中的所有内容，都是通过名为[空间](https://confluence.atlassian.com/conf68/spaces-947170008.html)的机制组织的。

这里的一个空间，就对应到本地的文件系统目录。

## 文档内容生成规则

首先，所有以`.`开头的文件或者文件夹、以及符号链接、设备等非普通的文件，都会被忽略。

除此之外，其他的文件/文件夹，按照如下规则处理：

### 文件夹

- 名为assets的文件夹及其所属的子文件夹，都不会生成任何页面。
- 其他名称的文件夹会生成一个同名的页面。
   - 如果文件夹内有名为`index`后缀为`.md`或`.xml`的索引文件时，索引文件内容会作为文件夹对应页面的内容。
   - 当文件夹内没有索引文件时，会填充缺省的内容。缺省内容是名为[Children Display](https://confluence.atlassian.com/doc/children-display-macro-139501.html)的Confluence宏，该宏自动替换为该页面的子页面索引。

### 普通文件


- **assets目录下的文件**：
	会被视作附件，上传到**所在的assets目录的父目录**对应的页面中。
- **其他目录下的文件**：

	- **以.md为后缀名**：会被当作Markdown内容解析。解析后的内容上传作为Confluence内容，内容的标题去掉后缀后的文件名部分。
	- **以.xml为后缀名**：会被当作原生的[Confluence Storage](https://confluence.atlassian.com/doc/confluence-storage-format-790796544.html)上传为Confluence的内容。内容标题仍然是去掉后缀名后的文件名部分。
	- **其他后缀名**：会被视作附件。会上传到其**所在的目录**对应的页面中，附件文件名/标题就是文件名。

## Markdown撰写指南

### 本地图片和附件链接

如果需要插入图片或者附件链接，建议将图片/附件文件放在Markdown文件同级目录的assets子目录中。并通过相对路径的方式来引用。例如，有如下目录结构：

```bash
tree ./
 |- Hallo.md
 `- assets/
     |- my_picture.jpg
     `- my_attachment.zip
```

那么，在*Hallo.md*文件中，我们就可以像下面这样来创建图片和附件的链接：

```markdown
![图片](assets/my_picture.jpg)
![下载链接](assets/my_attachment.zip)
```

当然，如果你喜欢，也可以把图片/附件放到Markdown相同的目录中。不过这种方式下可能管理起来不太方便。

不管放在同级目录，还是assets子目录里，我们都可以自动替换掉对应的链接地址为Confluence的附件地址。使得转换到Confluence中的效果，和本地编辑器的效果一致。


> ***提示***
>
> 作为目录索引文件的`index.md`中引用的图片/附件，请放在和目录的父目录中，**与目录本身保持同级**，而不是和`index.md`文件保持同级。


### 流程图、时序图、甘特图


目前很多Markdown编辑器，支持将流程图、时序图、甘特图的描述性语言，自动预览为对应的SVG图片。他们中的大多数都是通过一些第三方的Javascript库实现的。比如常见的[mermaid](https://mermaidjs.github.io/)。


经过和IT部门的同事合作，我们也内置了对于`mermaid`库的支持。


在本地编辑器（推荐[typora](https://typora.io/)）编辑时，你只需要像其他语言一样，撰写`mermaid`语言的代码片段。本地编辑器会自动渲染为预览的流程图。


当Markdown文本同步到Confluence后，我们也会自动渲染为同样的SVG图片。


例如，你可以像下面这样撰写一个简单的流程图

    ```mermaid
    graph TD
       Start --> Stop
    ```

### Confluence页面的链接

所有的Confluence页面，都支持display的方式进行链接。所以你可以通过下面这种方式创建页面间的链接：

```markdown
[标题](https://wiki.tap4fun.com/display/Confluence空间代码/页面标题#锚点名)
```

例如：

```markdown
[Chat](https://wiki.tap4fun.com/display/TGS/Chat)
```
