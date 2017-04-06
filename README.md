# AWS architecture

Here is the design of the final aws cloud solution to handle our different use cases in the cloud.

![AWS Architecture](https://github.com/scality/picsou/blob/master/docs/aws-architecture.png)

# Picsou

[![Circle CI](http://ci.ironmann.io/gh/scality/picsou.svg?style=svg)](http://ci.ironmann.io/gh/scality/picsou)

Picsou is a Cloud base costs reporter and manager.

## Build

#### Local

```
$ make
```

#### Docker

```
$ make docker-build
```

## Usage

#### Local

```
$ ./picsou report
```

#### Docker

```
$ docker build -t picsou .
$ docker run -v -e PICSOU_USER={{EmailAddress}} -e PICSOU_PSD={{EmailPassword}} $HOME/.aws:/root/.aws picsou
```

## Support

Please open an issue to receive support for this project.

## Contributing

Create a new branch, make your changes, and open a pull request.
