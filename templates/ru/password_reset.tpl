<html>
  <body>
    <p>Здравствуйте {{ .user.email }}!</p>
    <p>
      Используйте эту уникальную ссылку для сброса вашего пароля:
      <a href="{{ .domain }}/accounts/password_reset?reset_token={{ .token }}&lang=ru">Сбросить</a>
    </p>
  </body>
</html>
