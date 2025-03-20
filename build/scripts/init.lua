#!/usr/bin/tarantool

require "0-add-spaces"
require "1-add-privileges"
require "2-add-roles"
require "3-add-users"
require "4-add-role-privileges"
require "5-add-role-users"

box.cfg {
    listen = '0.0.0.0:33011',
    log_level = 4 -- warn
}

box.schema.user.passwd('admin', 'password')

--- init database
box.once('init', function()
    add_spaces()
    add_privileges()
    add_roles()
    add_users()
    add_role_privileges()
    add_role_users()
end)
