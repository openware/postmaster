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

| Variable            | Description                         | Required | Default              |
|---------------------|-------------------------------------|----------|----------------------|
| `JWT_PUBLIC_KEY`    | RSA Public Key for decoding         | *yes*    |                      |
| `SENDGRID_API_KEY`  | SendGrid API Key                    | *yes*    |                      |
| `RABBITMQ_HOST`     | Host of RabbitMQ daemon             | *no*     | `localhost`          |
| `RABBITMQ_PORT`     | Port of RabbitMQ daemon             | *no*     | `5672`               |
| `RABBITMQ_USERNAME` | RabbitMQ username                   | *no*     | `guest`              |
| `RABBITMQ_PASSWORD` | RabbitMQ password                   | *no*     | `guest`              |
| `SENDER_EMAIL`      | Email address of mail sender        | *no*     | `noreply@postmaster.com` |
| `SENDER_NAME `      | Name of mail sender                 | *no*     | `postmaster`             |
| `CONFIRM_URL`       | URL template for confirmation email | *no*     | `http://www.example.com/accounts/confirmation?confirmation_token=#{}` |

## License

Project released under the terms of the MIT [license](./LICENSE).
