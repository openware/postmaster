From: {{ .From }}
To: {{ .To }}
Subject: Confirmation Instructions

<html>
    <body>
        <p>Hello {{ .Record.Email }}</p>
        <p>
            <a href="{{ .Record.ConfirmationUri }}">Confirm your account</a>
        </p>
    </body>
</html>