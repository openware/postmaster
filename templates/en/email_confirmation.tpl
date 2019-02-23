<html>
  <body>
    <p>Hello {{ .user.email }}!</p>
    <p>
      Use this unique link to confirm your password:
      <a href="http://www.example.com/{{ .token }}">Confirm</a>
    </p>
  </body>
</html>
