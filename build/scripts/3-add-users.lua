#!/usr/bin/tarantool

function add_users()
    local s = box.space.user
    -- The password used during database initialization
    -- This needs to be changed for further work
    local pass = '$2a$10$IXCXaTHCMZOs6u85ArI7y.jsuti79y4Rjmp6OFcoAv/l23amn9VYe' -- 12345
    --
    s:insert{nil, 'Администратор', 'admin', pass, false}
end
