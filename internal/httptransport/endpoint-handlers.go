package httptransport

// func(f strictecho.StrictEchoHandlerFunc) strictecho.StrictEchoHandlerFunc {
// 		extractTokenValue := func(header http.Header) (string, error) {
// 			values, ok := header["Authorization"]
// 			if !ok {
// 				return "", fmt.Errorf("token not found")
// 			}

// 			if len(values) < 1 {
// 				return "", fmt.Errorf("token not found")
// 			}

// 			ul := strings.Split(values[0], " ")

// 			if len(ul) < 2 { //nolint:mnd
// 				return "", fmt.Errorf("invalid token")
// 			}

// 			return ul[1], nil
// 		}

// 		return func(f strictecho.StrictEchoHandlerFunc) strictecho.StrictEchoHandlerFunc {
// 			return func(ctx echo.Context, request any) (any, error) {
// 				value, err := extractTokenValue(ctx.Request().Header)
// 				if err != nil {
// 					return nil, err
// 				}

// 				token, err := authidjwt.ParseToken([]byte(signingKey), []byte(value))
// 				if err != nil {
// 					return nil, err //nolint:wrapcheck
// 				}

// 				if !token.Valid {
// 					return nil, fmt.Errorf("token not valid")
// 				}

// 				if err := searchPrivilegeFunc(ctx.Request().Context(), token.SessionID, code); err != nil {
// 					return nil, fmt.Errorf("searchPrivilegeFunc | %w", err)
// 				}

// 				return f(ctx, request)
// 			}
// 		}
// 	}
// )
