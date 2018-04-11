package ui

import ()

const header = `
<!DOCTYPE html>
<html>
<head>
	<title>Secure File Share</title>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width,initial-scale=1.0,maximum-scale=1.0,user-scalable=no" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/semantic-ui/2.3.1/semantic.css">
</head>
<body>
`

const footer = `
</body>
</html>
`

// PasswordPage will usk user to enter the passcode
const PasswordPage = header + `
<div class="ui container">
<div class="ui info message">
  <div class="header">
    Attention!
  </div>
  If you failed to enter the passcode 3 times, the server will start self-destruction <br>
  You have {{.}} chances left.
</div>
<form class="ui form" method="post">
	<div class="field">
	    <label>Passcode</label>
	    <input type="password" name="code" placeholder="Please Enter Your Passcode">
  </div>
  <button class="ui blue button" type="submit">Submit</button>
</form>
</div>` + footer

// DownloadPage shows a list of items
const DownloadPage = header + `
<div class="ui container">
<div class="ui info message">
  <div class="header">
    Attention!
  </div>
  Please click destroy, after download all the files
</div>
	<table>
		<thead>
			<tr>
				<th>File</th>
				<th>Action</th>
			</tr>
		</thead>
		<tbody>
			{{range .}}
		    <tr>
			    <td>{{.}}</td>
			    <td><a href="/download/{{.}}">download</a></td>
		    </tr>
		    {{end}}
		</tbody>    
		<tfoot>
			<tr>
				<th colspan="2"><a class="ui red button" href="/destroy/">destroy</a></th>
			</tr>
		</tfoot>
	</table>
</div>
` + footer
