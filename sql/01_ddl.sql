-- DDL
create table if not exists TaskGroup (
    id integer primary key not null,
    group_name text not null,
    sys_created datetime not null default current_timestamp,
    sys_updated datetime not null default current_timestamp,
    sys_status integer not null default 1
);

create table if not exists TaskStatus(
    id integer primary key not null,
    status_name text not null,
    sys_created datetime not null default current_timestamp,
    sys_updated datetime not null default current_timestamp,
    sys_status integer not null default 1
);

create table if not exists Task (
    id integer primary key not null,
    summary text not null,
    priority integer not null default 500,
    date_due datetime null,
    date_completed datetime null,
    description text null,
    closing_comment text null,
    status_id integer not null default 1,
    group_id integer not null default 1,
    parent_id integer null,
    sys_created datetime not null default current_timestamp,
    sys_updated datetime not null default current_timestamp,
    sys_status integer not null default 1,
    
    foreign key (status_id) references TaskStatus(id),
    foreign key (group_id) references TaskGroup(id),
    foreign key (parent_id) references Task(id)
);

create table if not exists AuditLog (
    id integer primary key not null,
    event_date datetime not null default current_timestamp,
    operation text not null,
    target_table text not null,
    target_id integer not null,
    target_column text not null,
    new_value text null
);

create table if not exists SysVersion (
    module text primary key not null,
    version text not null
);
