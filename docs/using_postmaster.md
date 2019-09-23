# Using Postmaster

This guide explains the basics of using postmaster to manage events from rubykube Event API.

Nowadays postmaster support two components, that produce events

* [Barong](https://github.com/rubykube/barong)
* [Peatio](https://github.com/rubykube/peatio)

Read more about **peatio** Event API [here](https://github.com/rubykube/peatio/blob/master/docs/api/event_api.md).

Read more about **barong** Event API [here](https://github.com/rubykube/barong/blob/master/docs/event_api.md).

## Concepts

An *Event* is a message produced to message broker in [RFC7515](https://tools.ietf.org/html/rfc7515).

Usually, events have the following structure

```JSON
{
  "payload": "string",
  "signatures": [
    {
      "header": {
        "kid":"string"
      },
      "protected":"string",
      "signature":"string"
    }
  ]
}
```

An *Expression* gives you the ability to apply additional event matching.

## Deep dive

By specification defined in *RFC7515* we can build a JSON Web Token.

After parsing of prebuild JWT, you will receive this payload with this structure.

```JSON
{
  "iss": "string",
  "jti": "string",
  "iat": 1567777420,
  "exp": 1567777480,
  "event": {
    "record": {
      "user": {
        "uid": "string",
        "email": "string",
      },
      "language": "string",
    },
    "name": "string"
  }
}
```

**Note:** All events should be properly signed using RS256 private key.

Postmaster validates the signature of each upcoming event.

Further, we will see that postmaster user can work with this payload directly.

```JSON
{
  "record": {
    "user": {
      "uid": "UID12345678",
      "email": "johndoe@example.com"
    },
    "language": "EN"
  },
  "changes": {},
  "name": "string"
}
```

## Expressions

When core model updates, the provider sends an event with routing key `<model>.updated`, but sometimes you need
to apply additional event matching, in this case, you may want to use `expression`.

The expression is evaluated by the postmaster when it receives an event.

In expression there two fields you can access: `record` and `changes`.

**Note:** Your `expression` should always return boolean!

Example:

```yaml
events:
- name: Withdrawal Succeed
  key: withdraw.updated
  exchange: peatio
  expression: changes.state in ["errored", "confirming"] && record.state == "succeed"
  templates:
    EN:
      subject: Withdrawal Succeed
      template_path: templates/en/withdraw_succeed.tpl

```

For details about **experssion** language read [documentation](https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md).

## Configuration

All languages, that will be used in the application, should be defined using `languages` field.

```yaml
languages:
- code: EN
  name: English
- code: RU
  name: Russian
```

Each Event API provider uses own AMQP exchange and algorithm to sign payload.

```yaml
exchanges:
  barong:
    name: barong.events.system
    signer: peatio
  peatio:
    name: peatio.events.model
    signer: peatio
```

Using keychain algorithms and defined public keys for each provider postmaster will validate the data.

```yaml
keychain:
  barong:
    algorithm: RS256
    value: "public_key"
  peatio:
    algorithm: RS256
    value: "public_key"
```

In `events` you may define any event type from Event API providers and prepare email template for it.

```yaml
events:
- name: Email Confirmation
  key: user.email.confirmation.token
  exchange: barong
  templates:
    EN:
      subject: Registration Confirmation
      template: Hello, {{ .record.user.email }}!
    RU:
      subject: Подтверждение Регистрации
      template_path: templates/ru/email_confirmation.tpl
```

For `templates` you can use either mount templates directly to container and *template_path* or use *template*.

## Installation

### Locally

Build the binary

```sh
go build ./cmd/postmaster/postmaster.go
```

Start your postmaster

```sh
./postmaster
```

Output:

```
XX:XXAM INF waiting for events
XX:XXAM INF successfully connected to amqp://xxxxx:xxxxx@localhost:5672/
```

### Docker

Use official docker image for postmaster quay.io/openware/postmaster, which is automatically built on each commit in the master branch.

```sh
docker run quay.io/openware/postmaster --name=postmaster
```
