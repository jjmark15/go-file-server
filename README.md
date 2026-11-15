# Go server

A zero-dependency, single-binary static file server written in Go.
Designed as a faster, more robust alternative to `python -m http.server`.

## Features

- **zero dependencies**: uses only go standard library.
- **single binary**: easy to distribute and deploy.
- **https support**: generates self-signed certificates in-memory on the fly.
- **logging**: simple request logging similar to python's format.
- **concurrent**: handles multiple connections efficiently via go routines.

## Usage

Run the server in the current directory on port 8000 by default:

``` bash
go-file-server
```

### Options

Plaintext
  -p string
        port to serve on (default "8000")
  -d string
        directory to serve (default ".")
  -h string
        host interface to bind to (default "0.0.0.0")
  -tls
        enable https with self-signed certificate
  -v
        show version

### Examples

Serve a specific directory on a custom port:

``` bash
go-file-server -p 9090 -d /var/www/html
```

Serve strictly on localhost (not accessible from network):

``` bash
go-file-server -h 127.0.0.1
```

Enable https mode (useful for testing secure contexts locally):

``` bash
go-file-server -tls -p 443
```

NOTE: since the certificate is self-signed and generated in-memory, your browser will warn you. This is expected.

## License

MIT
