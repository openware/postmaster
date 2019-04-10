<html>
  <body>
    <p>Здравствуйте {{ .user.email }}!</p>
    <p>
      Используйте эту уникальную ссылку для подтверждения вашего пароля:
      <a href="{{ .domain }}/accounts/confirmation?confirmation_token={{ .token }}&lang=ru">Подтвердить</a>
    </p>
  </body>
</html>
