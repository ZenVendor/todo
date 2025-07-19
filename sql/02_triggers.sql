-- TRIGGERS

-- Create Tasks
--drop trigger TRG_Audit_Create_Tasks_Main;
create trigger TRG_Audit_Create_Tasks_Main 
    after insert on Tasks
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'create', 'Tasks', new.id, 'summary', new.summary);
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'create', 'Tasks', new.id, 'priority', new.priority);
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'create', 'Tasks', new.id, 'group_id', new.group_id);
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'create', 'Tasks', new.id, 'status_id', new.status_id);
    end;

--drop trigger TRG_Audit_Create_Tasks_Due;
create trigger TRG_Audit_Create_Tasks_Due 
    after insert on Tasks
    when new.date_due is not null
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'create', 'Tasks', new.id, 'date_due', new.date_due);
    end;

--drop trigger TRG_Audit_Create_Tasks_Parent;
create trigger TRG_Audit_Create_Tasks_Parent 
    after insert on Tasks
    when new.parent_id is not null
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'create', 'Tasks', new.id, 'parent_id', new.parent_id);
    end;

--drop trigger TRG_Audit_Create_Tasks_Description;
create trigger TRG_Audit_Create_Tasks_Description 
    after insert on Tasks
    when new.description is not null
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'create', 'Tasks', new.id, 'description', new.description);
    end;

-- Update Tasks
--drop trigger TRG_Audit_Update_Tasks_summary;
create trigger TRG_Audit_Update_Tasks_summary
    after update of summary on Tasks
    when new.summary <> old.summary
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'update', 'Tasks', new.id, 'summary', new.summary);
    end;

--drop trigger TRG_Audit_Update_Tasks_priority;
create trigger TRG_Audit_Update_Tasks_priority
    after update of priority on Tasks
    when new.priority <> old.priority
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'update', 'Tasks', new.id, 'priority', new.priority);
    end;

--drop trigger TRG_Audit_Update_Tasks_date_due;
create trigger TRG_Audit_Update_Tasks_date_due
    after update of date_due on Tasks
    when new.date_due <> old.date_due
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'update', 'Tasks', new.id, 'date_due', new.date_due);
    end;

--drop trigger TRG_Audit_Update_Tasks_date_completed;
create trigger TRG_Audit_Update_Tasks_date_completed
    after update of date_completed on Tasks
    when new.date_completed <> old.date_completed
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'update', 'Tasks', new.id, 'date_completed', new.date_completed);
    end;

--drop trigger TRG_Audit_Update_Tasks_description;
create trigger TRG_Audit_Update_Tasks_description
    after update of description on Tasks
    when new.description <> old.description
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'update', 'Tasks', new.id, 'description', new.description);
    end;

--drop trigger TRG_Audit_Update_Tasks_closing_comment;
create trigger TRG_Audit_Update_Tasks_closing_comment
    after update of closing_comment on Tasks
    when new.closing_comment <> old.closing_comment
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'update', 'Tasks', new.id, 'closing_comment', new.closing_comment);
    end;

--drop trigger TRG_Audit_Update_Tasks_status_id;
create trigger TRG_Audit_Update_Tasks_status_id
    after update of status_id on Tasks
    when new.status_id <> old.status_id
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'update', 'Tasks', new.id, 'status_id', new.status_id);
    end;

--drop trigger TRG_Audit_Update_Tasks_group_id;
create trigger TRG_Audit_Update_Tasks_group_id
    after update of group_id on Tasks
    when new.group_id <> old.group_id
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'update', 'Tasks', new.id, 'group_id', new.group_id);
    end;

--drop trigger TRG_Audit_Update_Tasks_parent_id;
create trigger TRG_Audit_Update_Tasks_parent_id
    after update of parent_id on Tasks
    when new.parent_id <> old.parent_id
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'update', 'Tasks', new.id, 'parent_id', new.parent_id);
    end;

--drop trigger TRG_Audit_Soft_Delete_Tasks;
create trigger TRG_Audit_Soft_Delete_Tasks 
    after update of sys_status on Tasks
    when new.sys_status = 0
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'delete', 'Tasks', new.id, 'sys_status', new.sys_status);
    end;

--drop trigger TRG_Audit_Soft_Restore_Tasks;
create trigger TRG_Audit_Soft_Restore_Tasks 
    after update of sys_status on Tasks
    when new.sys_status = 1
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'restore', 'Tasks', new.id, 'sys_status', new.sys_status);
    end;

-- Delete Tasks
--drop trigger TRG_Audit_Hard_Delete_Tasks;
create trigger TRG_Audit_Hard_Delete_Tasks 
    after delete on Tasks
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'destroy', 'Tasks', old.id, 'summary', old.summary);
    end;


-- Create Groups
--drop trigger TRG_Audit_Create_Groups;
create trigger TRG_Audit_Create_Groups 
    after insert on Groups
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'create', 'Groups', new.id, 'group_name', new.group_name);
    end;

-- Update Groups
--drop trigger TRG_Audit_Update_Groups_name;
create trigger TRG_Audit_Update_Groups_name
    after update of name on Groups
    when new.group_name <> old.group_name
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'update', 'Groups', new.id, 'group_name', new.group_name);
    end;

--drop trigger TRG_Audit_Soft_Delete_Groups;
create trigger TRG_Audit_Soft_Delete_Groups 
    after update of sys_status on Groups
    when new.sys_status = 0
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'delete', 'Groups', new.id, 'sys_status', new.sys_status);
    end;

--drop trigger TRG_Audit_Soft_Restore_Groups;
create trigger TRG_Audit_Soft_Restore_Groups 
    after update of sys_status on Groups
    when new.sys_status = 1
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'restore', 'Groups', new.id, 'sys_status', new.sys_status);
    end;

-- Delete Groups
--drop trigger TRG_Audit_Hard_Delete_Groups;
create trigger TRG_Audit_Hard_Delete_Groups 
    after delete on Groups
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'destroy', 'Groups', old.id, 'group_name', old.group_name);
    end;


-- Create Statuses
--drop trigger TRG_Audit_Create_Statuses;
create trigger TRG_Audit_Create_Statuses 
    after insert on Statuses
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'create', 'Statuses', new.id, 'status_name', new.status_name);
    end;

-- Update Statuses
--drop trigger TRG_Audit_Update_Statuses_name;
create trigger TRG_Audit_Update_Statuses_name
    after update of name on Statuses
    when new.status_name <> old.status_name
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'update', 'Statuses', new.id, 'status_name', new.status_name);
    end;

--drop trigger TRG_Audit_Soft_Delete_Statuses;
create trigger TRG_Audit_Soft_Delete_Statuses 
    after update of sys_status on Statuses
    when new.sys_status = 0
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'delete', 'Statuses', new.id, 'sys_status', new.sys_status);
    end;

--drop trigger TRG_Audit_Soft_Restore_Groups;
create trigger TRG_Audit_Soft_Restore_Statuses 
    after update of sys_status on Statuses
    when new.sys_status = 1
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'restore', 'Statuses', new.id, 'sys_status', new.sys_status);
    end;

-- Delete Statuses
--drop trigger TRG_Audit_Hard_Delete_Statuses;
create trigger TRG_Audit_Hard_Delete_Statuses 
    after delete on Statuses
    begin
        insert into AuditLog (event_date, operation, target_table, target_id, target_column, new_value) 
        values (current_timestamp, 'destroy', 'Statuses', old.id, 'status_name', old.status_name);
    end;
