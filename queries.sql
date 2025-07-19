-- Table exists
select count(*) 
from sqlite_schema 
where 
    type = 'table' 
    and tbl_name in ('tasklist', 'taskgroup');

-- Create tables
create schema todo;
create table todo.task (
    id integer primary key not null,
    group_id integer not null default 1,
    status_id integer not null default 1,
    priority integer not null default 500,
    summary text not null,
    description text,
    closing_comment text,
    date_due datetime,
    date_completed datetime,
    parent_id integer,
    sys_created datetime not null default datetime(now),
    sys_updated datetime not null default datetime(now),
    sys_deleted int not null default 0
);
create table todo.group (
    id integer primary key not null,
    name text not null,
    sys_created datetime not null default datetime(now),
    sys_updated datetime not null default datetime(now),
    sys_deleted int not null default 0
);
create table todo.status (
    id integer primary key not null,
    name text not null,
    sys_created datetime not null default datetime(now),
    sys_updated datetime not null default datetime(now),
    sys_deleted int not null default 0
);

-- Sys inserts
insert into todo.group (id, name) values (1, 'Default');
insert into todo.status (id, name) values 
    (1, 'New'),
    (2, 'In Progress'),
    (3, 'On Hold'),
    (4, 'Completed');

-- Insert task
insert into todo.task (
    group_id
    , priority
    , summary
    , description
    , date_due
    , parent_id
) values (?, ?, ?, ?, ?, ?);

-- Insert group
insert into todo.group (name) values (?);

-- Select task
create view todo.taskslist_all as (
select 
    t.id
    , t.summary
    , t.priority
    , t.done
    , t.description
    , t.comment
    , t.due
    , t.completed
    , t.parent_id
    , t.created
    , t.updated
    , g.id
    , g.name
from 
    todo.task t
    join todo.group g on g.id = t.group_id
where 
    t.id = ?;

-- Select group by name
select 
    g.id
    , g.name
    , g.created
    , g.updated
from taskgroup g
where g.Name = ?;

-- Select group by Id
select 
    g.id
    , g.name
    , g.created
    , g.updated
from taskgroup g
where g.id = ?;

-- Update task
update tasklist 
set 
    summary = ?
    , priority = ?
    , group_id = ?
    , done = ?
    , description = ?
    , comment = ?
    , due = ?
    , completed = ?
    , task_id = ?
    , updated = datetime(now)
where id = ?;

-- Update group
update taskgroup 
set 
    name = ?
    , updated = datetime(now) 
where id = ?;

-- Delete task
delete from tasklist where id = ?;

-- Delete group
delete from taskgroup where id = ?;

-- count all, open, closed, overdue
select count(*) from tasklist;
select count(*) from tasklist where done = 0;
select count(*) from tasklist where done = 1;
select count(*) from tasklist where done = 0 and due < date(now);

-- Count in groups all, open, closed, overdue
select g.name, count(*) as cnt
from 
    taskgroup g
    join tasklist t on t.group_id = g.id
group by g.name
order by cnt;

select g.name, count(*) as cnt
from
    taskgroup g
    join tasklist t on t.group_id = g.id
where t.done = 0
group by g.name
order by cnt;

select g.name, count(*) as cnt
from
taskgroup g
join tasklist t on t.group_id = g.id
where t.done = 1
group by g.name
order by cnt;

select g.name, count(*) as cnt
from
    taskgroup g
    join tasklist t on t.group_id = g.id
where 
    t.done = 0 
    and t.due < date(now)
group by g.name
order by cnt;

-- List
select 
    t.id
	, t.summary
	, t.priority
	, t.done
    , t.description
    , t.comment
	, t.due
	, t.completed
    , t.task_id
	, t.created
	, t.updated
	, g.id
	, g.name
	, g.created
	, g.updated
from
    tasklist t
    join taskgroup g on g.id = t.group_id
order by
    t.done
	, t.priority desc
	, t.due
	, g.name; 

select 
    t.id
	, t.summary
	, t.priority
	, t.done
    , t.description
    , t.comment
	, t.due
	, t.completed
    , t.task_id
	, t.created
	, t.updated
	, g.id
	, g.name
	, g.created
	, g.updated
from
    tasklist t
    join taskgroup g on g.id = t.group_id
where
    t.done = 0
order by
    t.priority desc
	, t.due nulls last
	, g.name;

select 
    t.id
	, t.summary
	, t.priority
	, t.done
    , t.description
    , t.comment
	, t.due
	, t.completed
    , t.task_id
	, t.created
	, t.updated
	, g.id
	, g.name
	, g.created
	, g.updated
from
    tasklist t
    join taskgroup g on g.id = t.group_id
where
    t.done = 1
order by
    g.name
	, t.completed desc;


select 
    t.id
	, t.summary
	, t.priority
    , t.description
    , t.closing_comment
	, t.due
	, t.completed
    , t.task_id
	, t.created
	, t.updated
	, g.id
	, g.name
	, g.created
	, g.updated
from
    tasklist t
    join taskgroup g on g.id = t.group_id
where
    t.done = 0
    and t.due < ?
order by
    t.priority desc
	, t.due
	, g.name;
