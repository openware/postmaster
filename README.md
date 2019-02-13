# postmaster

> :incoming_envelope: Mail events consumer.

[Barong](https://www.github.com/rubykube/barong) Event API Client.

## Overview

Consume mail events from barong and send emails over SMTP.

![Overview](./resources/overview.png)

## Usage

Start worker by running command below

```sh
$ go run ./cmd/postmaster/main.go
```

### Environment variables

| Variable            | Description                          | Required | Default              |
|---------------------|--------------------------------------|----------|----------------------|
| `JWT_PUBLIC_KEY`    | RSA Public Key for decoding          | *yes*    |                      |
| `RABBITMQ_HOST`     | Host of RabbitMQ daemon              | *no*     | `localhost`          |
| `RABBITMQ_PORT`     | Port of RabbitMQ daemon              | *no*     | `5672`               |
| `RABBITMQ_USERNAME` | RabbitMQ username                    | *no*     | `guest`              |
| `RABBITMQ_PASSWORD` | RabbitMQ password                    | *no*     | `guest`              |
| `SMTP_PASSWORD`     | Password used for auth to SMTP       | *yes*    |                      |
| `SMTP_PORT`         | Post of SMTP server                  | *no*     | `25`                 |
| `SMTP_HOST`         | Host of SMTP server                  | *no*     | `smtp.sendgrid.net`  |
| `SMTP_USER`         | User used for auth to SMTP           | *no*     | `apikey`             |
| `SENDER_EMAIL`      | Email address of mail sender         | *yes*    |                      |
| `SENDER_NAME `      | Name of mail sender                  | *no*     | `Postmaster`                   |
| `CONFIRM_URL`       | URL template for confirmation email  | *no*     | `http://example.com/#{}`       |
| `RESET_URL`         | URL template for reset password      | *no*     | `http://example.com/#{}`       |
| `SIGN_UP_TEMPLATE_PATH`        | Path to file with content | *no*     | `templates/sign_up.tpl`        |
| `PASSWORD_RESET_TEMPLATE_PATH` | Path to file with content | *no*     | `templates/password_reset.tpl` |

### Custom Emails

There is an ability to create custom message per event.

Mount file to container and export environment variable with it's path.

Supported events:
- Sign Up - `SIGN_UP_TEMPLATE_PATH`
- Password Reset - `PASSWORD_RESET_TEMPLATE_PATH`

We are using `EVENT_NAME_TEMPLATE_PATH` convention for variables with template paths.

## License

Project released under the terms of the MIT [license](./LICENSE).
