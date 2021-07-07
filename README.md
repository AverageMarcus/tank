# Tank

A webhook handler to send messages to a Matrix homeserver.

## Example

```sh
curl -X "POST" "http://localhost:3000/" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "message": "Test webhook",
  "type": "text",
  "roomID": "example"
}'
```

## Features

* Send text messages to a room
* Send notice messages to a room

## Building from source

With Docker:

```sh
make docker-build
```

Standalone:

```sh
make build
```

## Contributing

If you find a bug or have an idea for a new feature please [raise an issue](issues/new) to discuss it.

Pull requests are welcomed but please try and follow similar code style as the rest of the project and ensure all tests and code checkers are passing.

Thank you 💛

## License

See [LICENSE](LICENSE)
