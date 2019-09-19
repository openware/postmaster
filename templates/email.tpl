From: "{{ .FromName }}" <{{ .FromAddress }}>
To: {{ .ToAddress }}
Subject: {{ .Subject }}
Content-type: text/html; charset="UTF-8"
<html>

<head>
    <style>
    {{ .Style }}
    </style>
</head>

<body>
    {{ .Body }}
</body>

</html>