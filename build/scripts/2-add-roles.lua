#!/usr/bin/tarantool

function add_roles()
    local s = box.space.role
    --
    s:insert{nil, 'admin', 'Администраторы казино', 'Какие-то администраторы', false}
end