-- Table exists
select count(*) 
from sqlite_schema 
where 
    type = 'table' 
    and tbl_name in ('tasklist', 'taskgroup');

-- Create tables
create table tasklist (
    id integer primary key not null,
    short text not null,
    priority integer not null default 50,
    group_id integer not null,
    done integer not null default 0,
    long text,
    comment text,
    due datetime,
    completed datetime,
    task_id integer,
    created datetime not null default datetime(now),
    updated datetime not null default datetime(now)
);
create table taskgroup (
    id integer primary key not null,
    name text not null,
    created datetime not null default datetime(now),
    updated datetime not null default datetime(now)
);
insert into taskgroup (name, created, updated) values ('Default', );

-- Insert task
insert into tasklist (
    short
    , priority
    , group_id
    , long
    , comment
    , due
    , task_id
) values (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- Insert group
insert into taskgroup (name) values (?);

-- Select task
select 
    t.id
    , t.short
    , t.priority
    , t.done
    , t.long
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
    short = ?
    , priority = ?
    , group_id = ?
    , done = ?
    , long = ?
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
	, t.short
	, t.priority
	, t.done
    , t.long
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
	, t.short
	, t.priority
	, t.done
    , t.long
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
	, t.short
	, t.priority
	, t.done
    , t.long
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
	, t.short
	, t.priority
	, t.done
    , t.long
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
    and t.due < ?
order by
    t.priority desc
	, t.due
	, g.name;
