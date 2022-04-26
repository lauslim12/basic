// Package basic provides plug and play, generic, secure, easy to use, customizable, and painless Basic Authentication for Go's HTTP handlers.
// This package tries its best to implement all specifications in a customizable way as specified in
// [RFC 7617](https://datatracker.ietf.org/doc/html/rfc7617), the newest version of Basic Authentication which obsoletes
// [RFC 2617](https://datatracker.ietf.org/doc/html/rfc2617).
//
// Basic Authentication itself is a simple and secure way to protect your API endpoints. However, for it to be completely secure,
// you have to augment the authentication by using SSL/TLS. You may use hashes / encryption in advance, but I think it's not necessary.
// SSL/TLS provides excellent security as long as you can trust your Certificate Authority and can ensure your connections are end-to-end
// encrypted, no sniffers or spoofers whatsoever.
//
// In order to use this package, you have to instantiate a new `BasicAuth` instance before hooking it to a handler.
// You can set various configurations, such as the authenticator function, charset (if using non-standard, other than UTF-8 charset),
// invalid credentials response, invalid scheme response, custom realm, and static users if applicable. If you want to do anything
// with the `*http.Request` struct, it is recommended for you to process it in a previous custom middleware before implementing your
// authentication with this library. This package tries its best to be as generic as possible, so you can definitely use any web framework or
// customized handlers as long as it conforms to the main interface (`http.Handler`).
//
// As a note about the `BasicAuth` attributes, you may use the authenticator function in order to perform a more
// sophisticated authentication logic, such as pulling your user based on their username from the database. Another thing to note is that
// you can pass `nil` or `make(map[string]string)` to the `Users` attribute if you do not need static credentials. Finally, the
// `WWW-Authenticate` header is only sent if both `Charset` and `Realm` are set. `Users` attribute is a 1-to-1 mapping of username
// and password.
//
// See example in `example/main.go`.
package basic

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"net/http"
)

// BasicAuth is used to configure all the library options.
type BasicAuth struct {
	Authenticator              func(username, password string) bool // Custom callback to find out the validity of a user's authentication process. This can be implemented in any implementation detail (for example: DB calls).
	Charset                    string                               // Custom charset to be passed in the `WWW-Authenticate` header. According to RFC 7617, this has to be 'UTF-8'.
	InvalidCredentialsResponse http.Handler                         // Callback to be invoked after receiving an InvalidCredentials error.
	InvalidSchemeResponse      http.Handler                         // Callback to be invoked after receiving an InvalidScheme error.
	Realm                      string                               // Specific realm for an authorization endpoint. This can be an arbitrary string.
	Users                      map[string]string                    // Static credentials for all users. Can be `nil` if need be.
}

// NewCustomBasicAuth is used to set up Basic Auth options with customizable configurations.
func NewCustomBasicAuth(
	authenticator func(username, password string) bool,
	charset string,
	invalidCredentialsResponse http.Handler,
	invalidSchemeResponse http.Handler,
	realm string,
	users map[string]string,
) *BasicAuth {
	// Populate parameters with default values for several necessary attributes.
	defaultConfig := NewDefaultBasicAuth(users)
	if authenticator == nil {
		authenticator = defaultConfig.Authenticator
	}

	if invalidCredentialsResponse == nil {
		invalidCredentialsResponse = defaultConfig.InvalidCredentialsResponse
	}

	if invalidSchemeResponse == nil {
		invalidSchemeResponse = defaultConfig.InvalidSchemeResponse
	}

	if len(users) == 0 {
		users = defaultConfig.Users
	}

	return &BasicAuth{
		Authenticator:              authenticator,
		Charset:                    charset,
		InvalidCredentialsResponse: invalidCredentialsResponse,
		InvalidSchemeResponse:      invalidSchemeResponse,
		Realm:                      realm,
		Users:                      users,
	}
}

// NewDefaultBasicAuth is used to set up Basic Auth options with default configurations.
func NewDefaultBasicAuth(users map[string]string) *BasicAuth {
	return &BasicAuth{
		// Accepts username and password. If the list of users is populated, the function will
		// check whether the username exists and then tries to securely compare the passwords. If the list of users
		// does not exist / has the length of zero, the function will return false.
		Authenticator: func(username, password string) bool {
			if len(users) != 0 {
				if val, ok := users[username]; ok {
					// Small trick to prevent timing attacks by hashing both usernames and passwords before comparing
					// them. This has its own overhead, but completely prevents timing attacks.
					usernamesMatch := CompareInputs(username, users[username])
					passwordsMatch := CompareInputs(password, val)

					return usernamesMatch && passwordsMatch
				}
			}

			return false
		},

		// RFC 7617: Only accept `UTF-8`.
		Charset: "UTF-8",

		// Response that will be sent if the credentials are invalid.
		InvalidCredentialsResponse: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Invalid username and/or password!", http.StatusUnauthorized)
		}),

		// Response that will be sent if the scheme (header) is invalid.
		InvalidSchemeResponse: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Invalid authentication scheme!", http.StatusUnauthorized)
		}),

		// Custom realm in the authentication process.
		Realm: "",

		// List of users allowed to access the endpoint.
		Users: users,
	}
}

// SendInvalidCredentialsResponse is used to send back an invalid response if the
// Basic Authorization credentials are invalid.
func (a *BasicAuth) SendInvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	a.SetWWWAuthenticate(w)
	a.InvalidCredentialsResponse.ServeHTTP(w, r)
}

// SendInvalidSchemeResponse is used to send back invalid response if the Basic
// Authorization header is not in the proper format.
func (a *BasicAuth) SendInvalidSchemeResponse(w http.ResponseWriter, r *http.Request) {
	a.SetWWWAuthenticate(w)
	a.InvalidSchemeResponse.ServeHTTP(w, r)
}

// SetWWWAuthenticate sets the `WWW-Authenticate` network header on the API response payload. If the
// charset and realm are both empty, we do not set the `WWW-Authenticate` header.
func (a *BasicAuth) SetWWWAuthenticate(w http.ResponseWriter) {
	if a.Realm != "" && a.Charset != "" {
		realm := fmt.Sprintf(`Basic realm="%s", charset="%s"`, a.Realm, a.Charset)
		w.Header().Set("WWW-Authenticate", realm)
	}
}

// Authenticate is a middleware to safeguard a route with the updated version of Basic
// Authentication (RFC 7617).
func (a *BasicAuth) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Grabs the username and password of the Basic Authentication.
		username, password, ok := r.BasicAuth()
		if !ok {
			a.SendInvalidSchemeResponse(w, r)
			return
		}

		// Try to authenticate the user.
		authenticated := a.Authenticator(username, password)

		// If not match, return 401.
		if !authenticated {
			a.SendInvalidCredentialsResponse(w, r)
			return
		}

		// If match, go to the next middleware.
		next.ServeHTTP(w, r)
	}
}

// CompareInputs is to safe compare two inputs (prevents timing attacks).
func CompareInputs(input, expected string) bool {
	// Hash input and expected with fast-hash.
	inputHash := sha256.Sum256([]byte(input))
	expectedHash := sha256.Sum256([]byte(expected))

	// Return boolean value with timing-safe comparisons to know whether the values
	// passed as arguments are equal or not.
	return (subtle.ConstantTimeCompare(inputHash[:], expectedHash[:]) == 1)
}
