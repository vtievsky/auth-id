package sessionsvc

import "errors"

var ErrInvalidAccessTokenTTL = errors.New("duration of the access token is less than the duration of the refresh token")
