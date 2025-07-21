-- TRIGGERS

-- Create Task
create trigger if not exists TRG_Audit_Create_Task_Main 
    after insert on Task
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'create', 'Task', new.id, 'summary', new.summary);
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'create', 'Task', new.id, 'priority', new.priority);
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'create', 'Task', new.id, 'group_id', new.group_id);
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'create', 'Task', new.id, 'status_id', new.status_id);
    end;

create trigger if not exists TRG_Audit_Create_Task_Due 
    after insert on Task
    when new.date_due is not null
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'create', 'Task', new.id, 'date_due', new.date_due);
    end;

create trigger if not exists TRG_Audit_Create_Task_Parent 
    after insert on Task
    when new.parent_id is not null
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'create', 'Task', new.id, 'parent_id', new.parent_id);
    end;

create trigger if not exists TRG_Audit_Create_Task_Description 
    after insert on Task
    when new.description is not null
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'create', 'Task', new.id, 'description', new.description);
    end;

-- Update Task
create trigger if not exists TRG_Audit_Update_Task_summary
    after update of summary on Task
    when new.summary <> old.summary
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'update', 'Task', new.id, 'summary', new.summary);
    end;

create trigger if not exists TRG_Audit_Update_Task_priority
    after update of priority on Task
    when new.priority <> old.priority
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'update', 'Task', new.id, 'priority', new.priority);
    end;

create trigger if not exists TRG_Audit_Update_Task_date_due
    after update of date_due on Task
    when new.date_due <> old.date_due
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'update', 'Task', new.id, 'date_due', new.date_due);
    end;

create trigger if not exists TRG_Audit_Update_Task_date_completed
    after update of date_completed on Task
    when new.date_completed <> old.date_completed
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'update', 'Task', new.id, 'date_completed', new.date_completed);
    end;

create trigger if not exists TRG_Audit_Update_Task_description
    after update of description on Task
    when new.description <> old.description
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'update', 'Task', new.id, 'description', new.description);
    end;

create trigger if not exists TRG_Audit_Update_Task_closing_comment
    after update of closing_comment on Task
    when new.closing_comment <> old.closing_comment
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'update', 'Task', new.id, 'closing_comment', new.closing_comment);
    end;

create trigger if not exists TRG_Audit_Update_Task_status_id
    after update of status_id on Task
    when new.status_id <> old.status_id
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'update', 'Task', new.id, 'status_id', new.status_id);
    end;

create trigger if not exists TRG_Audit_Update_Task_group_id
    after update of group_id on Task
    when new.group_id <> old.group_id
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'update', 'Task', new.id, 'group_id', new.group_id);
    end;

create trigger if not exists TRG_Audit_Update_Task_parent_id
    after update of parent_id on Task
    when new.parent_id <> old.parent_id
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'update', 'Task', new.id, 'parent_id', new.parent_id);
    end;

create trigger if not exists TRG_Audit_Soft_Delete_Task 
    after update of sys_status on Task
    when new.sys_status = 0
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'delete', 'Task', new.id, 'sys_status', new.sys_status);
    end;

create trigger if not exists TRG_Audit_Soft_Restore_Task 
    after update of sys_status on Task
    when new.sys_status = 1
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'restore', 'Task', new.id, 'sys_status', new.sys_status);
    end;

-- Delete Task
create trigger if not exists TRG_Audit_Hard_Delete_Task 
    after delete on Task
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'destroy', 'Task', old.id, 'summary', old.summary);
    end;


-- Create TaskGroup
create trigger if not exists TRG_Audit_Create_Group 
    after insert on TaskGroup
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'create', 'Group', new.id, 'group_name', new.group_name);
    end;

-- Update TaskGroup
create trigger if not exists TRG_Audit_Update_Group_name
    after update of name on TaskGroup
    when new.group_name <> old.group_name
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'update', 'Group', new.id, 'group_name', new.group_name);
    end;

create trigger if not exists TRG_Audit_Soft_Delete_Group 
    after update of sys_status on TaskGroup
    when new.sys_status = 0
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'delete', 'Group', new.id, 'sys_status', new.sys_status);
    end;

create trigger if not exists TRG_Audit_Soft_Restore_Group 
    after update of sys_status on TaskGroup
    when new.sys_status = 1
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'restore', 'Group', new.id, 'sys_status', new.sys_status);
    end;

-- Delete TaskGroup
create trigger if not exists TRG_Audit_Hard_Delete_Group 
    after delete on TaskGroup
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'destroy', 'Group', old.id, 'group_name', old.group_name);
    end;


-- Create TaskStatus
create trigger if not exists TRG_Audit_Create_Status 
    after insert on TaskStatus
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'create', 'Status', new.id, 'status_name', new.status_name);
    end;

-- Update TaskStatus
create trigger if not exists TRG_Audit_Update_Status_name
    after update of name on TaskStatus
    when new.status_name <> old.status_name
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'update', 'Status', new.id, 'status_name', new.status_name);
    end;

create trigger if not exists TRG_Audit_Soft_Delete_Status 
    after update of sys_status on TaskStatus
    when new.sys_status = 0
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'delete', 'Status', new.id, 'sys_status', new.sys_status);
    end;

create trigger if not exists TRG_Audit_Soft_Restore_Status 
    after update of sys_status on TaskStatus
    when new.sys_status = 1
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'restore', 'Status', new.id, 'sys_status', new.sys_status);
    end;

-- Delete TaskStatus
create trigger if not exists TRG_Audit_Hard_Delete_Status 
    after delete on TaskStatus
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'destroy', 'Status', old.id, 'status_name', old.status_name);
    end;
