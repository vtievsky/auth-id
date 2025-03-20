#!/usr/bin/tarantool

function add_privileges()
    local s = box.space.privilege
    --
    s:insert{nil, 'user_read', 'Чтение справочника пользователей', ''}
    s:insert{nil, 'user_create', 'Создание пользователя', ''}
    s:insert{nil, 'user_update', 'Изменение пользователя', ''}
    s:insert{nil, 'user_delete', 'Удаление пользователя', ''}
    --
    s:insert{nil, 'pass_reset', 'Сброс пароля', ''}
    s:insert{nil, 'pass_change', 'Изменение собственного пароля', ''}
    --
    s:insert{nil, 'role_read', 'Чтение справочника ролей', ''}
    s:insert{nil, 'role_create', 'Создание роли', ''}
    s:insert{nil, 'role_update', 'Изменение роли', ''}
    s:insert{nil, 'role_delete', 'Удаление роли', ''}
    --
    s:insert{nil, 'user2role_read', 'Чтение ролей пользователя', ''}
    s:insert{nil, 'user2privilege_read', 'Чтение привилегий пользователя', ''}
    s:insert{nil, 'user2session_read', 'Чтение сессий пользователя', ''}
    s:insert{nil, 'user2session_delete', 'Удаление сессии пользователя', ''}
    --
    s:insert{nil, 'role2user_read', 'Чтение пользователeй роли', ''}
    s:insert{nil, 'role2user_create', 'Добавление роли пользователю', ''}
    s:insert{nil, 'role2user_update', 'Изменение роли пользователя', ''}
    s:insert{nil, 'role2user_delete', 'Удаление роли пользователя', ''}
    --
    s:insert{nil, 'role2privilege_read', 'Чтение привилегий роли', ''}
    s:insert{nil, 'role2privilege_create', 'Добавление привилегии роли', ''}
    s:insert{nil, 'role2privilege_update', 'Изменение привилегии роли', ''}
    s:insert{nil, 'role2privilege_delete', 'Удаление привилегии роли', ''}
    --
    -- s:insert{nil, 'search_session_privilege', 'Проверка наличия привилегии у сессии', ''}
end