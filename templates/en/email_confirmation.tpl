<html>
  <body>
    <p>Hello {{ .user.email }}!</p>
    <p>
      Use this unique link to confirm your password:
      <a href="{{ .domain }}/accounts/confirmation?confirmation_token={{ .token }}&lang={{ .language }}">Confirm</a>
    </p>
  </body>
</html>
