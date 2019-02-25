<html>
  <body>
    <p>Здравствуйте {{ .user.email }}!</p>
    <p>
      Используйте эту уникальную ссылку для сброса вашего пароля:
      <a href="http://www.example.com/{{ .token }}">Сбросить</a>
    </p>
  </body>
</html>
