<html>
    <body>
        <p>Hello {{ .User.Email }}</p>
        <p>
            Use this unique link to reset your password: <a href="{{ .ResetPasswordURI }}">Reset</a>
        </p>
    </body>
</html>
