# Basic

Provides plug and play, generic, secure, easy to use, customizable, and painless Basic Authentication middleware for Go's HTTP handlers. No dependencies!

This package tries its best to implement all specifications in a customizable way as specified in [RFC 7617](https://datatracker.ietf.org/doc/html/rfc7617), the newest version of Basic Authentication which obsoletes [RFC 2617](https://datatracker.ietf.org/doc/html/rfc2617).

## Why Basic?

- **No dependencies.** Requires only needs standard Go.
- **Battle-tested.** This library conforms to the standard library (which a lot of people use nowadays).
- **Lightweight.** Basic is small in size, due to not having any dependencies.
- **Secure.** Tries its best to implement as many security considerations as possible, but you **definitely have to use HTTPS in production if you intend to use this in production**.
- **Generic.** This library is generic and implements `http.Handler` to ensure maximum compatibility with as many Go frameworks as possible.
- **100% tested.** As this library is small, the code coverage is still 100% for now.
- **Well documented.** Check out this `README.md` document and the technical documentation for further reading!

## Security Considerations

If you want to use this in production environment, here are additional security considerations:

- Ensure you are running this using HTTPS with SSL/TLS to prevent man in the middle attacks.
- Enable HSTS (`Strict-Transport-Security`) to prevent your site from being accessed with HTTP. Set redirects (`301 Moved Permanently`) from HTTP to HTTPS permanently in your reverse proxy / Go app. Use HTTPS forever!
- Use secure HTTP headers to prevent malicious browser agents (`X-XSS-Protection`, `X-Content-Type-Options`, `X-DNS-Prefetch-Control`, and the like).
- Use rate limiters in endpoints protected by Basic Authentication to prevent brute-force attacks.
- As usual, keep your passwords strong. Use symbols, numbers, uppercases, and lowercases. Even better if you use password managers.
- Follow and read security guidelines: [OWASP Cheatsheets](https://cheatsheetseries.owasp.org/)!
- My two cents and security tip: Basic Authentication should placed in an endpoint that gives out sessions / tokens on successful authentication. Make sure that endpoint is not cacheable (use `PUT`, `PATCH`, `POST` without `Cache-Control` headers, by default they are not cacheable, do not use `GET` and `HEAD` if possible). This relieves the pain of having to deal with logout and/or cache problems. You can then delegate your authentication via the given out sessions / tokens.

## Documentation

Complete documentation could be seen in the official [pkg.go.dev site](https://pkg.go.dev/github.com/lauslim12/basic).

## Installation

You have to perform the following steps (assume using Go 1.18):

- Download this library.

```bash
go install github.com/lauslim12/basic

# for older go versions: go get -u github.com/lauslim12/basic
```

- Import it in your source code.

```go
import "github.com/lauslim12/basic"
```

- Instantiate the `BasicAuth` object, and you can wrap it in any endpoint you desire to protect it!

```go
func main() {
    users := map[string]string{"nehemiah":"nehemiah"}

    // Use default authenticator function, set charset to UTF-8, use default invalid scheme response,
    // use default invalid credentials response, set custom realm, and set static user list.
    basicAuth := basic.NewCustomBasicAuth(nil, "UTF-8", nil, nil, "Private", users)
    http.HandleFunc("/", basicAuth.Authenticate(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(http.StatusText(http.StatusOK)))
    }))
}
```

- Test the endpoint!

```bash
curl -u nehemiah:nehemiah <API_ENDPOINT_URL>
```

- Done!

## Customizations

Customization is the core part of this library. You can customize anything, and you can even define / create a middleware before or after the `Authenticate` middleware method if you need to perform some preprocessing or postprocessing.

- As an example, you may define your own authorizer if you need to do so. Below code is for reference:

```go
func main() {
    // This pseudocode example sets no static users and calls the user from the DB based on
    // the user's input. It then matches the password and returns the boolean value.
    basicAuth := basic.NewCustomBasicAuth(func(username, password string) bool {
        user := getUserFromDB(username)
        match := basic.CompareInputs(password, user.Password)

        return match
    }, "UTF-8", nil, nil, "Private Not-Static", nil)

    // After defining it, we then hook it into our handler.
    http.HandleFunc("/", basicAuth.Authenticate(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusAccepted)
        w.Write([]byte(http.StatusText(http.StatusAccepted)))
    }))
}
```

- You can customize your `Authenticator` function (signature is `func(username, password string) bool`), `Charset` (defaults to `UTF-8` according to RFC 7617), `InvalidSchemeResponse` (signature is `http.Handler`), `InvalidCredentialsResponse` (signature is `http.Handler`), `Realm` (signature is `string`), and `Users` (signature is `map[string]string`). As long as it conforms to the interface / function signature, you can customize it with anything you want.

## Examples

Please see examples at [the example project (`example/main.go`)](./example). You can run it by doing `go run example/main.go` and then connect to `localhost:5000` on your web browser / API client.

## Contributing

This tool is open source and the contribution of this tool is highly encouraged! If you want to contribute to this project, please feel free to read the `CONTRIBUTING.md` file for the contributing guidelines.

## License

This work is licensed under MIT License. Please check the `LICENSE` file for more information.
