# Blossom Server

This Blossom server is developed for Zapstore to host applications executable and other relates blobs. It also supports whitelisting mechanism to only let developers to upload blobs.

# How to run?

You have to set environment variables defined in [the example file](./.env.example) first on the target system. Then you can build the project using:

```sh
go build .
```

> `make build` will do the same for you.

The you can run the blossom server using:

```sh
./blossom-server
```

# License

[MIT License](./LICENSE)
