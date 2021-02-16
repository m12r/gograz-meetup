# GoGraz-Meetup

This service abstracts the meetup.com API in order to fetch information like
the number of attendees for each meetup.

## How to Build

This project uses [mage][mage] as a build tool, so you only need [go][go]
installed on your system.

Just type the following in the project's root directory to build:

```shell
go run boostrap.go
```

To list the available targets, you can either run the following command or
have a look at the [magefile.go](magefile.go).

```shell
go run bootstrap.go -l
```

[go]: https://go.dev
[mage]: https://magefile.org

