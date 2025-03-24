package httptransport

// Соответствие endpoint и привилегии
var (
	//nolint:gochecknoglobals
	endpointWithout = map[string]struct{}{
		"post/v1/users/:login/sessions": {}, // Аутентификация пользователя
	}

	//nolint:gochecknoglobals
	endpointWithPrivileges = map[string]string{
		"put/v1/passresets/:login":  "pass_reset",
		"put/v1/passchanges/:login": "pass_change",
		// Пользователи
		"get/v1/users":           "user_read",
		"get/v1/users/:login":    "user_read",
		"post/v1/users":          "user_create",
		"put/v1/users/:login":    "user_update",
		"delete/v1/users/:login": "user_delete",
		// Роли и привилегии пользователя
		"get/v1/users/:login/roles":      "user2role_read",
		"get/v1/users/:login/privileges": "user2privilege_read",
		// Сессии пользователя
		"get/v1/users/:login/sessions":                "user2session_read",
		"delete/v1/users/:login/sessions/:session_id": "user2session_delete",
		// Роли
		"get/v1/roles/:code":    "role_read",
		"get/v1/roles":          "role_read",
		"post/v1/roles":         "role_create",
		"put/v1/roles":          "role_update",
		"delete/v1/roles/:code": "role_delete",
		// Пользователи роли
		"get/v1/roles/:code/users":           "role2user_read",
		"post/v1/roles/:code/users/:login":   "role2user_create",
		"put/v1/roles/:code/users/:login":    "role2user_update",
		"delete/v1/roles/:code/users/:login": "role2user_delete",
		// Привилегии роли
		"get/v1/roles/:code/privileges":                    "role2privilege_read",
		"post/v1/roles/:code/privileges/:privilege_code":   "role2privilege_create",
		"put/v1/roles/:code/privileges/:privilege_code":    "role2privilege_update",
		"delete/v1/roles/:code/privileges/:privilege_code": "role2privilege_delete",
		// Справочник привилегий
		"get/v1/privileges": "privilege_read",
	}
)
