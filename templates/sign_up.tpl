<html>
    <body>
        <p>Hello {{ .User.Email }}</p>
        <p>
            <a href="{{ .EmailConfirmationURI }}">Confirm your account</a>
        </p>
    </body>
</html>
