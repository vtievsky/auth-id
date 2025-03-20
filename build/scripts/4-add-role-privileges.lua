#!/usr/bin/tarantool

function add_role_privileges()
    local roles = box.space.role
    local privileges = box.space.privilege
    local role_privileges = box.space.role_privilege
    --
    local role_id = roles.index.secondary:select({'admin'})[1].id
    --
    for _, tuple in privileges:pairs() do
        role_privileges:insert{role_id, tuple[1], true}
    end
end
