--drop table Tasks;
create table Tasks (
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
    sys_status integer not null default 1
);

--drop table Groups;
create table Groups (
    id integer primary key not null,
    group_name text not null,
    sys_created datetime not null default current_timestamp,
    sys_updated datetime not null default current_timestamp,
    sys_status integer not null default 1
);

--drop table Statuses;
create table Statuses(
    id integer primary key not null,
    status_name text not null,
    sys_created datetime not null default current_timestamp,
    sys_updated datetime not null default current_timestamp,
    sys_status integer not null default 1
);

--drop table AuditLog;
create table AuditLog (
    id integer primary key not null,
    event_date datetime not null default current_timestamp,
    operation text not null,
    target_table text not null,
    target_id integer not null,
    target_column text not null,
    new_value text null
);


