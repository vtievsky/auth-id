#!/usr/bin/tarantool

function add_users()
    local s = box.space.user
    --
    s:insert{nil, 'Администратор', 'admin', '*****', false}
end
