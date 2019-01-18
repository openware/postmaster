<html>
    <body>
        <p>Hello {{ .User.Email }}</p>
        <p>
            <a href="{{ .ConfirmationURI }}">Confirm your account</a>
        </p>
    </body>
</html>