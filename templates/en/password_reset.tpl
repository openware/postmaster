<html>
  <body>
    <p>Hello {{ .user.email }}!</p>
    <p>
      Use this unique link to reset your password:
      <a href="http://www.example.com/{{ .token }}">Reset</a>
    </p>
  </body>
</html>
