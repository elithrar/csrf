package csrf

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

// MaxAge sets the maximum age (in seconds) of a CSRF token's underlying cookie.
// Defaults to 12 hours.
func MaxAge(age int) func(*csrf) error {
	return func(cs *csrf) error {
		cs.opts.MaxAge = age
		return nil
	}
}

// Domain sets the cookie domain. Defaults to the current domain of the request
// only (recommended).
//
// This should be a hostname and not a URL. If set, the domain is treated as
// being prefixed with a '.' - e.g. "example.com" becomes ".example.com" and
// matches "www.example.com" and "secure.example.com".
func Domain(domain string) func(*csrf) error {
	return func(cs *csrf) error {
		cs.opts.Domain = domain
		return nil
	}
}

// Path sets the cookie path. Defaults to the path the cookie was issued from
// (recommended).
//
// This instructs clients to only respond with cookie for that path and its
// subpaths - i.e. a cookie issued from "/register" would be included in requests
// to "/register/step2" and "/register/submit".
func Path(p string) func(*csrf) error {
	return func(cs *csrf) error {
		cs.opts.Path = p
		return nil
	}
}

// Secure sets the 'Secure' flag on the cookie. Defaults to true (recommended).
func Secure(s bool) func(*csrf) error {
	return func(cs *csrf) error {
		cs.opts.Secure = s
		return nil
	}
}

// HttpOnly sets the 'HttpOnly' flag on the cookie. Defaults to true (recommended).
func HttpOnly(h bool) func(*csrf) error {
	return func(cs *csrf) error {
		// Note that the function and field names match the case of the
		// related http.Cookie field instead of the "correct" HTTPOnly name
		// that golint suggests.
		cs.opts.HttpOnly = h
		return nil
	}
}

// ErrorHandler allows you to change the handler called when CSRF request
// processing encounters an invalid token or request. A typical use would be to
// provide a handler that returns a static HTML file with a HTTP 403 status. By
// default a HTTP 404 status and a plain text CSRF failure reason are served.
//
// Note that a custom error handler can also access the csrf.Failure(c, r)
// function to retrieve the CSRF validation reason from Goji's request context.
func ErrorHandler(h web.Handler) func(*csrf) error {
	return func(cs *csrf) error {
		cs.opts.ErrorHandler = h
		return nil
	}
}

// RequestHeader allows you to change the request header the CSRF middleware
// inspects. The default is X-CSRF-Token.
func RequestHeader(header string) func(*csrf) error {
	return func(cs *csrf) error {
		cs.opts.RequestHeader = header
		return nil
	}
}

// FieldName allows you to change the name value of the hidden <input> field
// generated by csrf.FormField. The default is {{ .csrfToken }}
func FieldName(name string) func(*csrf) error {
	return func(cs *csrf) error {
		cs.opts.FieldName = name
		return nil
	}
}

// setStore sets the store used by the CSRF middleware.
// Note: this is private (for now) to allow for internal API changes.
func setStore(s store) func(*csrf) error {
	return func(cs *csrf) error {
		cs.st = s
		return nil
	}
}

// parseOptions parses the supplied options functions and returns a configured
// csrf handler.
func parseOptions(h http.Handler, opts ...func(*csrf) error) *csrf {
	// Set the handler to call after processing.
	cs := &csrf{
		h: h,
	}

	// Default to true. See Secure & HttpOnly function comments for rationale.
	// Set here to allow package users to override the default.
	cs.opts.Secure = true
	cs.opts.HttpOnly = true

	// Range over each options function and apply it
	// to our csrf type to configure it. Options functions are
	// applied in order, with any conflicting options overriding
	// earlier calls.
	for _, option := range opts {
		option(cs)
	}

	return cs
}
