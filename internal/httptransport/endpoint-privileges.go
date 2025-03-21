package httptransport

// Соответствие operation_id и привилегии
func EndpointPrivilegesMiddlewareFuncs(
	signingKey string,
	sessionService SessionService,
) map[string]EndpointPrivilegesMiddlewareFunc {
	return map[string]EndpointPrivilegesMiddlewareFunc{
		// Авторизация
		"Login":      without(),
		"ResetPass":  withPrivilege(signingKey, sessionService.Search, "pass_reset"),
		"ChangePass": withPrivilege(signingKey, sessionService.Search, "pass_change"),
		// Пользователи
		"GetUser":    withPrivilege(signingKey, sessionService.Search, "user_read"),
		"GetUsers":   withPrivilege(signingKey, sessionService.Search, "user_read"),
		"CreateUser": withPrivilege(signingKey, sessionService.Search, "user_create"),
		"UpdateUser": withPrivilege(signingKey, sessionService.Search, "user_update"),
		"DeleteUser": withPrivilege(signingKey, sessionService.Search, "user_delete"),
		// Роли и привилегии пользователя
		"GetUserRoles":      withPrivilege(signingKey, sessionService.Search, "user2role_read"),
		"GetUserPrivileges": withPrivilege(signingKey, sessionService.Search, "user2privilege_read"),
		// Сессии пользователя
		"GetUserSessions":   withPrivilege(signingKey, sessionService.Search, "user2session_read"),
		"DeleteUserSession": withPrivilege(signingKey, sessionService.Search, "user2session_delete"),
		// Роли
		"GetRole":    withPrivilege(signingKey, sessionService.Search, "role_read"),
		"GetRoles":   withPrivilege(signingKey, sessionService.Search, "role_read"),
		"CreateRole": withPrivilege(signingKey, sessionService.Search, "role_create"),
		"UpdateRole": withPrivilege(signingKey, sessionService.Search, "role_update"),
		"DeleteRole": withPrivilege(signingKey, sessionService.Search, "role_delete"),
		// Пользователи роли
		"GetRoleUsers":   withPrivilege(signingKey, sessionService.Search, "role2user_read"),
		"AddRoleUser":    withPrivilege(signingKey, sessionService.Search, "role2user_create"),
		"UpdateRoleUser": withPrivilege(signingKey, sessionService.Search, "role2user_update"),
		"DeleteRoleUser": withPrivilege(signingKey, sessionService.Search, "role2user_delete"),
		// Привилегии роли
		"GetRolePrivileges":   withPrivilege(signingKey, sessionService.Search, "role2privilege_read"),
		"AddRolePrivilege":    withPrivilege(signingKey, sessionService.Search, "role2privilege_create"),
		"UpdateRolePrivilege": withPrivilege(signingKey, sessionService.Search, "role2privilege_update"),
		"DeleteRolePrivilege": withPrivilege(signingKey, sessionService.Search, "role2privilege_delete"),
	}
}
