package website

import (
	"fmt"
	"net/http"

	"github.com/iangudger/basicGoAPI/database"
)

func Handler(res http.ResponseWriter, req *http.Request, dbconn database.DB) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(res, `<!doctype html>
<html>
	<head>
		<meta charset=utf-8>
		<title>Basic Go API</title>
	</head>
	<body>
		<h1>Basic Go API</h1>
		<p>Open source basic API written in Go. Includes registration, login/logout and password changes/resets.</p>
		<p>Easy deployment to Heroku.</p>
	</body>
</html>
`)
}
