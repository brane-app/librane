package middleware

import (
	"github.com/brane-app/database-library"

	"context"
	"net/http"
	"strings"
)

/**
 * Reject unauthed users
 */
func MustAuth(request *http.Request) (modified *http.Request, ok bool, code int, r_map map[string]interface{}, err error) {
	code = 401
	var bearer string = strings.TrimPrefix(request.Header.Get("Authorization"), BEARER_PREFIX)

	var owner string
	if owner, ok, err = database.ReadTokenStat(bearer); err != nil || !ok {
		r_map = map[string]interface{}{"error": "bad_auth"}
	}

	modified = request.WithContext(context.WithValue(
		request.Context(),
		"requester",
		owner,
	))

	return
}

/**
 * Reject banned requests
 * required before: MustAuth to get the requester
 */
func RejectBanned(request *http.Request) (_ *http.Request, ok bool, code int, r_map map[string]interface{}, err error) {
	var owner string
	var owned bool
	if owner, owned = request.Context().Value("requester").(string); !owned {
		ok = true
		return
	}

	ok = false
	var banned bool
	if banned, err = database.IsBanned(owner); err != nil || banned {
		code = 403
		r_map = map[string]interface{}{"error": "banned"}
		return
	}

	ok = !banned
	return
}

/**
 * Reject users who are not moderators
 */
func MustModerator(request *http.Request) (_ *http.Request, ok bool, code int, r_map map[string]interface{}, err error) {
	code = 403

	var owner string
	var owned bool
	if owner, owned = request.Context().Value("requester").(string); !owned {
		return
	}

	ok, err = database.IsModerator(owner)
	return
}

/**
 * Reject users who are not admins
 */
func MustAdmin(request *http.Request) (_ *http.Request, ok bool, code int, r_map map[string]interface{}, err error) {
	code = 403

	var owner string
	var owned bool
	if owner, owned = request.Context().Value("requester").(string); !owned {
		return
	}

	ok, err = database.IsAdmin(owner)
	return
}
