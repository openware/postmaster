<html>
  <body>
    <p>Здравствуйте {{ .user.email }}!</p>
    <p>
      Используйте эту уникальную ссылку для подтверждения вашего пароля:
      <a href="http://www.example.com/{{ .token }}">Подтвердить</a>
    </p>
  </body>
</html>
