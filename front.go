package main

func getFrontTemplate() string {
	// Helper function to compile the template into the binary
	return `<!DOCTYPE html><html lang="en">
<head>
<meta charset="utf-8" />
<title>start</title>
<style>
html {
background: url({{.ImageURL}}) no-repeat center center fixed;
-webkit-background-size: cover;
-moz-background-size: cover;
-o-background-size: cover;
background-size: cover;
}
#reference {
position: fixed;
bottom: 0;
right: 0;
background-color: black;
color: #fff;
font-family: Helvetica, Arial, sans-serif;
display: inline;
padding: 0.5rem;

-webkit-box-decoration-break: clone;
box-decoration-break: clone;
}
a:link {color: #ffffff; text-decoration: underline; }
a:active {color: #ffffff; text-decoration: underline; }
a:visited {color: #ffffff; text-decoration: underline; }
a:hover {color: #ffffff; text-decoration: none; }
</style>
</head>
<body>
<p id="reference">Photo by <a href="https://unsplash.com/@{{.Username}}">{{.Name}}</a>  on <a href="https://unsplash.com">Unsplash</a></p>
</body>
</html>`
}
