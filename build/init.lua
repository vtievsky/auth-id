#!/usr/bin/tarantool

box.cfg {
    listen = '0.0.0.0:33011',
    log_level = 4 -- warn
}

box.schema.user.passwd('admin', 'password')

--- init database
box.once('init', function()
    -- users
    if not box.space.user then
        if not box.sequence.user_seq then
            box.schema.sequence.create('user_seq', {
                min = 1,
                start = 1
            })
        end
        ---
        local s = box.schema.space.create('user')
        --
        s:format({{
            name = 'id',
            type = 'unsigned'
        }, {
            name = 'name',
            type = 'string'
        }, {
            name = 'login',
            type = 'string'
        }, {
            name = 'password',
            type = 'string'
        }, {
            name = 'blocked',
            type = 'boolean'
        }})
        --
        s:create_index('pk', {
            sequence = 'user_seq',
            type = 'tree',
            parts = {'id'}
        })
        s:create_index('secondary', {
            type = 'tree',
            unique = true,
            parts = {'login'}
        })
    end

    -- roles
    if not box.space.role then
        if not box.sequence.role_seq then
            box.schema.sequence.create('role_seq', {
                min = 1,
                start = 1
            })
        end
        --
        local s = box.schema.space.create('role')
        --
        s:format({{
            name = 'id',
            type = 'unsigned'
        }, {
            name = 'code',
            type = 'string'
        }, {
            name = 'name',
            type = 'string'
        }, {
            name = 'description',
            type = 'string'
        }, {
            name = 'blocked',
            type = 'boolean'
        }})
        --
        s:create_index('pk', {
            sequence = 'role_seq',
            type = 'tree',
            parts = {'id'}
        })
        s:create_index('secondary', {
            type = 'tree',
            unique = true,
            parts = {'code'}
        })
    end

    -- privileges
    if not box.space.privilege then
        if not box.sequence.privilege_seq then
            box.schema.sequence.create('privilege_seq', {
                min = 1,
                start = 1
            })
        end
        --
        local s = box.schema.space.create('privilege')
        --
        s:format({{
            name = 'id',
            type = 'unsigned'
        }, {
            name = 'code',
            type = 'string'
        }, {
            name = 'name',
            type = 'string'
        }, {
            name = 'description',
            type = 'string'
        }})
        --
        s:create_index('pk', {
            sequence = 'privilege_seq',
            type = 'tree',
            parts = {'id'}
        })
        s:create_index('secondary', {
            type = 'tree',
            unique = true,
            parts = {'code'}
        })
    end

    -- role-users
    if not box.space.role_user then
        local s = box.schema.space.create('role_user')
        --
        s:format({{
            name = 'role_id',
            type = 'unsigned'
        }, {
            name = 'user_id',
            type = 'unsigned'
        }, {
            name = 'date_in',
            type = 'unsigned'
        }, {
            name = 'date_out',
            type = 'unsigned'
        }})
        --
        s:create_index('pk', {
            type = 'tree',
            parts = {'role_id', 'user_id'}
        })
        s:create_index('primary', {
            type = 'tree',
            unique = false,
            parts = {'role_id'}
        })
        s:create_index('secondary', {
            type = 'tree',
            unique = false,
            parts = {'user_id'}
        })
    end

    -- role-privileges
    if not box.space.role_privilege then
        local s = box.schema.space.create('role_privilege')
        --
        s:format({{
            name = 'role_id',
            type = 'unsigned'
        }, {
            name = 'privilege_id',
            type = 'unsigned'
        }, {
            name = 'allowed',
            type = 'boolean'
        }})
        --
        s:create_index('pk', {
            type = 'tree',
            parts = {'role_id', 'privilege_id'}
        })
        s:create_index('primary', {
            type = 'tree',
            unique = false,
            parts = {'role_id'}
        })
        s:create_index('secondary', {
            type = 'tree',
            unique = false,
            parts = {'privilege_id'}
        })
    end
end)

--- 
box.once('init-1', function()
    local s = box.space.privilege
    --
    s:insert{nil, 'user_read', 'Чтение справочника пользователей', ''}
    s:insert{nil, 'user_create', 'Создание пользователя', ''}
    s:insert{nil, 'user_update', 'Изменение пользователя', ''}
    s:insert{nil, 'user_delete', 'Удаление пользователя', ''}
    --
    s:insert{nil, 'role_read', 'Чтение справочника ролей', ''}
    s:insert{nil, 'role_create', 'Создание роли', ''}
    s:insert{nil, 'role_update', 'Изменение роли', ''}
    s:insert{nil, 'role_delete', 'Удаление роли', ''}
    --
    s:insert{nil, 'role2privilege_read', 'Чтение привилегий роли', ''}
    s:insert{nil, 'role2privilege_create', 'Добавление привилегии роли', ''}
    s:insert{nil, 'role2privilege_update', 'Изменение привилегии роли', ''}
    s:insert{nil, 'role2privilege_delete', 'Удаление привилегии роли', ''}
    --
    s:insert{nil, 'role2user_read', 'Чтение ролей пользователя', ''}
    s:insert{nil, 'role2user_create', 'Добавление роли пользователю', ''}
    s:insert{nil, 'role2user_update', 'Изменение роли пользователя', ''}
    s:insert{nil, 'role2user_delete', 'Удаление роли пользователя', ''}
end)

-- Роли
box.once('init-2', function()
    local s = box.space.role
    --
    s:insert{nil, 'admin', 'Администраторы казино', '', false}
end)

-- Пользователи
box.once('init-3', function()
    local s = box.space.user
    --
    s:insert{nil, 'Администратор', 'admin', '12345', false}
end)
