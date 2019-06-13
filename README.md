# postmaster

[![Build Status](https://ci.microkube.com/api/badges/openware/postmaster/status.svg)](https://ci.microkube.com/openware/postmaster)

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

| Variable               | Description                          | Required | Default              |
|------------------------|--------------------------------------|----------|----------------------|
| `POSTMASTER_ENV`       | Environment, reacts on "production"  | *no*     |                      |
| `POSTMASTER_LOG_LEVEL` | Level of logging                     | *no*     | `debug`              |
| `JWT_PUBLIC_KEY`       | RSA Public Key for decoding          | *yes*    |                      |
| `RABBITMQ_HOST`        | Host of RabbitMQ daemon              | *no*     | `localhost`          |
| `RABBITMQ_PORT`        | Port of RabbitMQ daemon              | *no*     | `5672`               |
| `RABBITMQ_USERNAME`    | RabbitMQ username                    | *no*     | `guest`              |
| `RABBITMQ_PASSWORD`    | RabbitMQ password                    | *no*     | `guest`              |
| `SMTP_PASSWORD`        | Password used for auth to SMTP       | *yes*    |                      |
| `SMTP_PORT`            | Post of SMTP server                  | *no*     | `25`                 |
| `SMTP_HOST`            | Host of SMTP server                  | *no*     | `smtp.sendgrid.net`  |
| `SMTP_USER`            | User used for auth to SMTP           | *no*     | `apikey`             |
| `SENDER_EMAIL`         | Email address of mail sender         | *yes*    |                      |
| `SENDER_NAME `         | Name of mail sender                  | *no*     | `Postmaster`         |

## License

Project released under the terms of the MIT [license](./LICENSE).
