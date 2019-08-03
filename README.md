# Linden Honey GoLang Scrapper

> Web scrapper for the lyrics of Yegor Letov - GoLang Edition

## Technologies

* [Colly](https://github.com/gocolly/colly)

## Usage

### Local

Start application:

```sh
go run main.go
```

Install missing dependencies:

```sh
make install
```

Start in development mode:

```sh
make start
```

Run in docker container:

```sh
docker-compose up -d
```

Go lint: 

```bash
go get golang.org/x/lint/golint
golint ./...
```

### Rest API

* `/` hello
* `/songs` letov songs

### Docker

Bootstrap db using docker-compose:

```sh
docker-compose up
```

Stop and remove containers, networks, images, and volumes:

```sh
docker-compose down
```
