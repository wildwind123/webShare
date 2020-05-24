package html

const HtmlTemplate = "<!DOCTYPE html>" +
	"<html lang=\"en\">" +
	"<head>" +
	"<meta charset=\"UTF-8\">" +
	"<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">" +

	"<title>WebShared</title>" +
	"</head>" +
	"<body>" +
	"<div>" +
	"<header>" +
	"<div id=\"header-container\">" +
	"{{.Header}} <br>" +
	"{{range .Folders}}" +
	"<a href=\"?path={{ .LinkFolder }}\">{{ .FolderName }}</a>" +
	"{{end}}" +
	"</div>" +
	"</header>" +
	"<main>" +
	"<div id=\"main-container\">" +
	"<div class=\"main-row\">" +
	"<div class=\"row\">" +
	"<div class=\"column1\">Name</div>" +
	"<div class=\"column2\">Type</div>" +
	"<div class=\"column3\">Size</div>" +
	"<div class=\"column4\">Modified</div>" +
	"</div>" +
	"</div>" +
	"{{range .Files}}" +
	"<div class=\"main-row\">" +
	"<div class=\"row\">" +
	"<div class=\"column1\"><a href=\"{{ .Link }}\">{{ .Name }}</a></div>" +
	"<div class=\"column2\">{{ .Type }}</div>" +
	"<div class=\"column3\">{{ .Size }}</div>" +
	"<div class=\"column4\">{{ .Modified }}</div>" +
	"</div>" +
	"</div>" +
	"{{end}}" +
	"</div>" +
	"</main>" +
	"</div>" +
	"<style>" +
	".row{" +
	"display: flex;" +
	"flex-wrap: wrap;" +
	"}" +
	".column1{" +
	"width: 300px;" +
	"word-wrap: break-word;" +
	"}" +
	".column2,.column3{" +
	"margin: 3px;" +
	"width: 100px;" +
	"word-wrap: break-word;" +
	"}" +
	".column4{" +
	"word-wrap: break-word;" +
	"width: 200px;" +
	"}" +

	"body{" +
	"display: flex;" +
	"justify-content: center;" +
	"}" +
	"#main-container{" +
	"display: flex;" +
	"flex-direction: column;" +
	"min-width: 200px;" +
	"overflow:auto;" +
	"padding-left: 10px;" +
	"padding-right: 10px;" +
	"padding-top: 5px;" +
	"}" +
	"header{" +
	"border-top-left-radius: 10px;" +
	"border-top-right-radius: 10px;" +
	"height: 50px;" +
	"background-color: #d5ebb9;" +
	"}" +
	"main{" +
	"background-color: #e9f2da;" +
	"border-bottom-left-radius: 10px;" +
	"border-bottom-right-radius: 10px;" +
	"box-shadow: #c1c9b7 0.5em 0.5em 0.3em;" +
	"}" +
	"#header-container{" +
	"padding-left: 10px;" +
	"padding-right: 10px;" +
	"padding-top: 5px;" +
	"}" +
	".size,.modified,.type{" +
	"text-align: center;" +
	"}" +
	".main-row{" +
	"background-color: #e2ebd3;" +
	"margin: 5px;" +
	"}" +

	"@media screen and (max-width: 1000px) {" +
	".column1,.column2,.column3,.column4 {" +
	"width: 100px;" +
	"}" +
	".main-row{" +
	"background-color: #c8d1b9;" +
	"}" +
	"}" +
	"</style>" +
	"</body>" +
	"</html>"
