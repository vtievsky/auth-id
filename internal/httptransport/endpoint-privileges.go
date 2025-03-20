package httptransport

// Соответствие operation_id и привилегии
func EndpointPrivilegesMiddlewareFuncs(
	signingKey string,
	sessionService SessionService,
) map[string]EndpointPrivilegesMiddlewareFunc {
	return map[string]EndpointPrivilegesMiddlewareFunc{
		// Авторизация
		"Login":      without(),
		"ResetPass":  withPrivilege(signingKey, sessionService.Find, "pass_reset"),
		"ChangePass": withPrivilege(signingKey, sessionService.Find, "pass_change"),
		// Пользователи
		"GetUser":    withPrivilege(signingKey, sessionService.Find, "user_read"),
		"GetUsers":   withPrivilege(signingKey, sessionService.Find, "user_read"),
		"CreateUser": withPrivilege(signingKey, sessionService.Find, "user_create"),
		"UpdateUser": withPrivilege(signingKey, sessionService.Find, "user_update"),
		"DeleteUser": withPrivilege(signingKey, sessionService.Find, "user_delete"),
		// Роли и привилегии пользователя
		"GetUserRoles":      withPrivilege(signingKey, sessionService.Find, "user2role_read"),
		"GetUserPrivileges": withPrivilege(signingKey, sessionService.Find, "user2privilege_read"),
		// Сессии пользователя
		"GetUserSessions":   withPrivilege(signingKey, sessionService.Find, "user2session_read"),
		"DeleteUserSession": withPrivilege(signingKey, sessionService.Find, "user2session_delete"),
		// Роли
		"GetRole":    withPrivilege(signingKey, sessionService.Find, "role_read"),
		"GetRoles":   withPrivilege(signingKey, sessionService.Find, "role_read"),
		"CreateRole": withPrivilege(signingKey, sessionService.Find, "role_create"),
		"UpdateRole": withPrivilege(signingKey, sessionService.Find, "role_update"),
		"DeleteRole": withPrivilege(signingKey, sessionService.Find, "role_delete"),
		// Пользователи роли
		"GetRoleUsers":   withPrivilege(signingKey, sessionService.Find, "role2user_read"),
		"AddRoleUser":    withPrivilege(signingKey, sessionService.Find, "role2user_create"),
		"UpdateRoleUser": withPrivilege(signingKey, sessionService.Find, "role2user_update"),
		"DeleteRoleUser": withPrivilege(signingKey, sessionService.Find, "role2user_delete"),
		// Привилегии роли
		"GetRolePrivileges":   withPrivilege(signingKey, sessionService.Find, "role2privilege_read"),
		"AddRolePrivilege":    withPrivilege(signingKey, sessionService.Find, "role2privilege_create"),
		"UpdateRolePrivilege": withPrivilege(signingKey, sessionService.Find, "role2privilege_update"),
		"DeleteRolePrivilege": withPrivilege(signingKey, sessionService.Find, "role2privilege_delete"),
	}
}
