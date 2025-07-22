-- all
create view if not exists task_list_all as
select 
    t.id 
    , t.summary
    , t.priority
    , t.date_due
    , t.date_completed
    , t.description
    , t.closing_comment
    , t.status_id
    , s.status_name
    , t.group_id
    , g.group_name
    , t.parent_id
    , t.sys_created
    , t.sys_updated
    , t.sys_status
from
    Task t
    inner join TaskGroup g on g.id = t.group_id
    inner join TaskStatus s on s.id = t.status_id
where 1=1
    and t.sys_status = 1;

-- new
create view if not exists task_list_new as
select 
    t.id 
    , t.summary
    , t.priority
    , t.date_due
    , t.date_completed
    , t.description
    , t.closing_comment
    , t.status_id
    , s.status_name
    , t.group_id
    , g.group_name
    , t.parent_id
    , t.sys_created
    , t.sys_updated
    , t.sys_status
from
    Task t
    inner join TaskGroup g on g.id = t.group_id
    inner join TaskStatus s on s.id = t.status_id
where 1=1
    and t.sys_status = 1
    and t.status_id = 1;

-- open
create view if not exists task_list_ongoing as
select 
    t.id 
    , t.summary
    , t.priority
    , t.date_due
    , t.date_completed
    , t.description
    , t.closing_comment
    , t.status_id
    , s.status_name
    , t.group_id
    , g.group_name
    , t.parent_id
    , t.sys_created
    , t.sys_updated
    , t.sys_status
from
    Task t
    inner join TaskGroup g on g.id = t.group_id
    inner join TaskStatus s on s.id = t.status_id
where 1=1
    and t.sys_status = 1
    and t.status_id in (1, 2, 3);

-- in_progress
create view if not exists task_list_in_progress as
select 
    t.id 
    , t.summary
    , t.priority
    , t.date_due
    , t.date_completed
    , t.description
    , t.closing_comment
    , t.status_id
    , s.status_name
    , t.group_id
    , g.group_name
    , t.parent_id
    , t.sys_created
    , t.sys_updated
    , t.sys_status
from
    Task t
    inner join TaskGroup g on g.id = t.group_id
    inner join TaskStatus s on s.id = t.status_id
where 1=1
    and t.sys_status = 1
    and t.status_id = 2;

-- on hold
create view if not exists task_list_on_hold as
select 
    t.id 
    , t.summary
    , t.priority
    , t.date_due
    , t.date_completed
    , t.description
    , t.closing_comment
    , t.status_id
    , s.status_name
    , t.group_id
    , g.group_name
    , t.parent_id
    , t.sys_created
    , t.sys_updated
    , t.sys_status
from
    Task t
    inner join TaskGroup g on g.id = t.group_id
    inner join TaskStatus s on s.id = t.status_id
where 1=1
    and t.sys_status = 1
    and t.status_id = 3;

-- completed
create view if not exists task_list_completed as
select 
    t.id 
    , t.summary
    , t.priority
    , t.date_due
    , t.date_completed
    , t.description
    , t.closing_comment
    , t.status_id
    , s.status_name
    , t.group_id
    , g.group_name
    , t.parent_id
    , t.sys_created
    , t.sys_updated
    , t.sys_status
from
    Task t
    inner join TaskGroup g on g.id = t.group_id
    inner join TaskStatus s on s.id = t.status_id
where 1=1
    and t.sys_status = 1
    and t.status_id = 4;

-- deleted
create view if not exists task_list_deleted as
select 
    t.id 
    , t.summary
    , t.priority
    , t.date_due
    , t.date_completed
    , t.description
    , t.closing_comment
    , t.status_id
    , s.status_name
    , t.group_id
    , g.group_name
    , t.parent_id
    , t.sys_created
    , t.sys_updated
    , t.sys_status
from
    Task t
    inner join TaskGroup g on g.id = t.group_id
    inner join TaskStatus s on s.id = t.status_id
where 1=1
    and t.sys_status = 0;

-- groups
create view if not exists group_list as 
select id, group_name from TaskGroup where sys_status = 1;

-- status
create view if not exists status_list as 
select id, status_name from TaskStatus where sys_status = 1;

-- counts
create view if not exists task_counts as 
select 
    count(*) as count_all
    , sum(case when status_id = 1 then 1 else 0 end) as count_new
    , sum(case when status_id = 2 then 1 else 0 end) as count_in_progress
    , sum(case when status_id = 3 then 1 else 0 end) as count_on_hold
    , sum(case when status_id = 4 then 1 else 0 end) as count_completed
    , sum(case when status_id in (1, 2, 3) then 1 else 0 end) as count_ongoing
from Task
where sys_status = 1;
