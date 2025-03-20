#!/usr/bin/tarantool

function add_role_users()
    local roles = box.space.role
    local users = box.space.user
    local role_users = box.space.role_user
    --
    local role_id = roles.index.secondary:select({'admin'})[1].id
    local user_id = users.index.secondary:select({'admin'})[1].id
    --
    local date_in = os.time({
        year = 2000,
        month = 1,
        day = 1
    })
    local date_out = os.time({
        year = 2999,
        month = 12,
        day = 31
    })
    --
    role_users:insert{role_id, user_id, date_in, date_out}
end
