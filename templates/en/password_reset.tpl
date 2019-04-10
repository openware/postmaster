<html>
  <body>
    <p>Hello {{ .user.email }}!</p>
    <p>
      Use this unique link to reset your password:
      <a href="{{ .domain }}/accounts/password_reset?reset_token={{ .token }}&lang={{ .language }}">Reset</a>
    </p>
  </body>
</html>
