package loghtml

const tpl = `
<!DOCTYPE html>
<html lang="es">

<head>
    <meta charset="UTF-8">

    <link rel="stylesheet" href="css/style.css">
</head>

<body>
    <div class='contetFilter'>
        <div id='contentOrder' class='sort' onclick="functionOrderAndClass()">
            Ordenar por Fecha
        </div>
        <div class='filter select'>
            <select id='typeData'>
                <option selected value="all">All</option>
                <option value="information">Info</option>
                <option value="debug">Debug</option>
                <option value="warning">Warn</option>
                <option value="error">Error</option>
            </select>
        </div>
    </div>

    <div class='content' id="content-list"></div>
    <script type="text/javascript" src="js/myScript.js"></script>
	{{range .Files}}<script type="text/javascript" src="js/{{ . }}"></script>{{end}}
</body>
</html>
`
