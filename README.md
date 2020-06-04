# Go-Libnss
This is an abstracted library in Go that wraps around libnss, allowing you to write new nss modules purely in Go.

## What's Working?
Right now we only have binding for `passwd`, `group`, and `shadow`. This may change in the future to also implement other NSS features. These are just the most obvious first targets.

## How do I use it?
First you should `go get` the package:
```
go get github.com/izanagi1995/go-libnss
```
Take a gander at the `example` directory. It includes information on how to compile the example and you can use the implementation as a boiler plate for your own project.

## License
MIT License

# Projects Using Go-Libnss
[Go-Libnss-Etcd](https://github.com/izanagi1995/go-libnss-etcd) - An NSS package that stores passwd, group, and shadow data in etcd. It also comes with CLI tools for managing users and groups stored in etcd.
