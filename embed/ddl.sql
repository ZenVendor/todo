-- DDL
create table if not exists Project (
    id integer primary key not null,
    project_name text not null,
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
    description text not null default '',
    closing_comment text not null default '',
    status integer not null default 1,
    project_id integer not null default 1,
    parent_id integer not null default 0,
    sys_created datetime not null default current_timestamp,
    sys_updated datetime not null default current_timestamp,
    sys_status integer not null default 1,
    
    foreign key (project_id) references Project(id),
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
    id integer primary key not null,
    cs_db text not null
);
