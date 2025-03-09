#!/usr/bin/tarantool

box.cfg {
    listen = '0.0.0.0:33011',
    log_level = 4 -- warn
}

box.schema.user.passwd('admin', 'password')

--- Init database
box.once('init', function()
    -- users
    if not box.space.user then
        if not box.sequence.user_seq then
            box.schema.sequence.create('user_seq', { min = 1, start = 1 })
        end
        ---
        local s = box.schema.space.create('user')
        --
        s:format({
            { name = 'id',   type = 'unsigned' },
            { name = 'name', type = 'string' },
            { name = 'login', type = 'string' },
            { name = 'blocked', type = 'boolean' },
        })
        s:create_index('pk', { sequence = 'user_seq', type = 'tree', parts = { 'id' } })
        s:create_index('secondary', { type = 'tree', parts = { 'login' } })
        --
        s:insert { nil, 'Пупкин Василий Иванович', 'pupkin_vi', false }
        s:insert { nil, 'Папироскина Мария Ивановна', 'papiroskina_mi', false }
    end

    -- roles
    if not box.space.role then
        if not box.sequence.role_seq then
            box.schema.sequence.create('role_seq', { min = 1, start = 1 })
        end
        --
        local s = box.schema.space.create('role')
        --
        s:format({
            { name = 'id',   type = 'unsigned' },
            { name = 'code', type = 'string' },
            { name = 'name', type = 'string' },
            { name = 'blocked', type = 'boolean' },
        })
        s:create_index('pk', { sequence = 'role_seq', type = 'tree', parts = { 'id' } })
        s:create_index('secondary', { type = 'tree', parts = { 'code' } })
        --
        s:insert { nil, 'admin', 'Администраторы казино', '', false }
    end

    -- -- Task states
    -- if not box.space.task_state then
    --     if not box.sequence.task_state_seq then
    --         box.schema.sequence.create('task_state_seq', { min = 1, start = 1 })
    --     end

    --     --
    --     local s = box.schema.space.create('task_state')
    --     s:format({
    --         { name = 'id',   type = 'unsigned' },
    --         { name = 'code', type = 'string' },
    --         { name = 'name', type = 'string' },
    --         { name = 'note', type = 'string' },
    --     })
    --     s:create_index('pk', { sequence = 'task_state_seq', type = 'tree', parts = { 'id' } })
    --     s:create_index('secondary', { type = 'tree', parts = { 'code' } })

    --     --
    --     s:insert { nil, 'busy', 'В процессе обработки', '' }
    --     s:insert { nil, 'error', 'Завершено с ошибкой', '' }
    --     s:insert { nil, 'success', 'Успешно завершено', '' }
    --     s:insert { nil, 'waiting', 'Ожидание обработки', '' }
    --     s:insert { nil, 'cancelled', 'Отменено', '' }
    -- end

    -- -- Task marks
    -- if not box.space.mark then
    --     if not box.sequence.mark_seq then
    --         box.schema.sequence.create('mark_seq', { min = 1, start = 1 })
    --     end

    --     --
    --     local s = box.schema.space.create('mark')
    --     s:format({
    --         { name = 'id',   type = 'unsigned' },
    --         { name = 'code', type = 'string' },
    --         { name = 'name', type = 'string' },
    --         { name = 'note', type = 'string' },
    --     })
    --     s:create_index('pk', { sequence = 'mark_seq', type = 'tree', parts = { 'id' } })
    --     s:create_index('secondary', { type = 'tree', parts = { 'code' } })
    -- end

    -- -- Task bool attributes
    -- if not box.space.task_attr_bool then
    --     if not box.sequence.task_attr_bool_seq then
    --         box.schema.sequence.create('task_attr_bool_seq', { min = 1, start = 1 })
    --     end

    --     --
    --     local s = box.schema.space.create('task_attr_bool')
    --     s:format({
    --         { name = 'id',      type = 'unsigned' },
    --         { name = 'task_id', type = 'unsigned' },
    --         { name = 'attr_id', type = 'unsigned' },
    --         { name = 'value',   type = 'boolean' },
    --         { name = 'note',    type = 'string' },
    --     })
    --     s:create_index('pk', { sequence = 'task_attr_bool_seq', type = 'tree', parts = { 'id' } })
    --     s:create_index('secondary', { type = 'tree', unique = false, parts = { 'task_id' } })
    -- end

    -- -- Task date attributes
    -- if not box.space.task_attr_date then
    --     if not box.sequence.task_attr_date_seq then
    --         box.schema.sequence.create('task_attr_date_seq', { min = 1, start = 1 })
    --     end

    --     --
    --     local s = box.schema.space.create('task_attr_date')
    --     s:format({
    --         { name = 'id',      type = 'unsigned' },
    --         { name = 'task_id', type = 'unsigned' },
    --         { name = 'attr_id', type = 'unsigned' },
    --         { name = 'value',   type = 'unsigned' },
    --         { name = 'note',    type = 'string' },
    --     })
    --     s:create_index('pk', { sequence = 'task_attr_date_seq', type = 'tree', parts = { 'id' } })
    --     s:create_index('secondary', { type = 'tree', unique = false, parts = { 'task_id' } })
    -- end

    -- -- Task integer attributes
    -- if not box.space.task_attr_int then
    --     if not box.sequence.task_attr_int_seq then
    --         box.schema.sequence.create('task_attr_int_seq', { min = 1, start = 1 })
    --     end

    --     --
    --     local s = box.schema.space.create('task_attr_int')
    --     s:format({
    --         { name = 'id',      type = 'unsigned' },
    --         { name = 'task_id', type = 'unsigned' },
    --         { name = 'attr_id', type = 'unsigned' },
    --         { name = 'value',   type = 'unsigned' },
    --         { name = 'note',    type = 'string' },
    --     })
    --     s:create_index('pk', { sequence = 'task_attr_int_seq', type = 'tree', parts = { 'id' } })
    --     s:create_index('secondary', { type = 'tree', unique = false, parts = { 'task_id' } })
    -- end

    -- -- Task numeric attributes
    -- if not box.space.task_attr_num then
    --     if not box.sequence.task_attr_num_seq then
    --         box.schema.sequence.create('task_attr_num_seq', { min = 1, start = 1 })
    --     end

    --     --
    --     local s = box.schema.space.create('task_attr_num')
    --     s:format({
    --         { name = 'id',      type = 'unsigned' },
    --         { name = 'task_id', type = 'unsigned' },
    --         { name = 'attr_id', type = 'unsigned' },
    --         { name = 'value',   type = 'number' },
    --         { name = 'note',    type = 'string' },
    --     })
    --     s:create_index('pk', { sequence = 'task_attr_num_seq', type = 'tree', parts = { 'id' } })
    --     s:create_index('secondary', { type = 'tree', unique = false, parts = { 'task_id' } })
    -- end

    -- -- Task text attributes
    -- if not box.space.task_attr_text then
    --     if not box.sequence.task_attr_text_seq then
    --         box.schema.sequence.create('task_attr_text_seq', { min = 1, start = 1 })
    --     end

    --     --
    --     local s = box.schema.space.create('task_attr_text')
    --     s:format({
    --         { name = 'id',      type = 'unsigned' },
    --         { name = 'task_id', type = 'unsigned' },
    --         { name = 'attr_id', type = 'unsigned' },
    --         { name = 'value',   type = 'string' },
    --         { name = 'note',    type = 'string' },
    --     })
    --     s:create_index('pk', { sequence = 'task_attr_text_seq', type = 'tree', parts = { 'id' } })
    --     s:create_index('secondary', { type = 'tree', unique = false, parts = { 'task_id' } })
    -- end

    -- -- Tasks
    -- if not box.space.task then
    --     if not box.sequence.task_seq then
    --         box.schema.sequence.create('task_seq', { min = 1, start = 1 })
    --     end

    --     --
    --     local s = box.schema.space.create('task')
    --     s:format({
    --         { name = 'id',            type = 'unsigned' },
    --         { name = 'id_task_type',  type = 'unsigned' },
    --         { name = 'id_task_state', type = 'unsigned' },
    --         { name = 'created_at',    type = 'unsigned' },
    --         { name = 'updated_at',    type = 'unsigned' },
    --         { name = 'started_at',    type = 'unsigned' },
    --         { name = 'delay',         type = 'unsigned' },
    --         { name = 'code',          type = 'string' },
    --         { name = 'note',          type = 'string' },
    --     })
    --     s:create_index('pk', { sequence = 'task_seq', type = 'tree', parts = { 'id' } })
    --     s:create_index('secondary', { type = 'tree', unique = false, parts = { 'id_task_state' } })
    --     s:create_index('updatedary', { type = 'tree', unique = false, parts = { 'updated_at' } })
    -- end

    -- -- Task marks
    -- if not box.space.task_mark then
    --     if not box.sequence.task_mark_seq then
    --         box.schema.sequence.create('task_mark_seq', { min = 1, start = 1 })
    --     end

    --     --
    --     local s = box.schema.space.create('task_mark')
    --     s:format({
    --         { name = 'id',      type = 'unsigned' },
    --         { name = 'task_id', type = 'unsigned' },
    --         { name = 'mark_id', type = 'unsigned' },
    --     })
    --     s:create_index('pk', { sequence = 'task_mark_seq', type = 'tree', parts = { 'id' } })
    --     s:create_index('secondary', { type = 'tree', unique = false, parts = { 'task_id' } })
    --     s:create_index('markery', { type = 'tree', unique = false, parts = { 'mark_id' } })
    -- end
end)
