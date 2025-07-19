-- all
create view task_list_all as
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
    Tasks t
    inner join Groups g on g.id = t.group_id
    inner join Statuses s on s.id = t.status_id
where 1=1
    and t.sys_status = 1
    and g.sys_status = 1;

-- new
create view task_list_new as
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
    Tasks t
    inner join Groups g on g.id = t.group_id
    inner join Statuses s on s.id = t.status_id
where 1=1
    and t.sys_status = 1
    and g.sys_status = 1
    and t.status_id = 1;

-- open
create view task_list_ongoing as
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
    Tasks t
    inner join Groups g on g.id = t.group_id
    inner join Statuses s on s.id = t.status_id
where 1=1
    and t.sys_status = 1
    and g.sys_status = 1
    and t.status_id in (1, 2, 3);

-- in_progress
create view task_list_in_progress as
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
    Tasks t
    inner join Groups g on g.id = t.group_id
    inner join Statuses s on s.id = t.status_id
where 1=1
    and t.sys_status = 1
    and g.sys_status = 1
    and t.status_id = 2;

-- on hold
create view task_list_on_hold as
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
    Tasks t
    inner join Groups g on g.id = t.group_id
    inner join Statuses s on s.id = t.status_id
where 1=1
    and t.sys_status = 1
    and g.sys_status = 1
    and t.status_id = 3;

-- completed
create view task_list_completed as
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
    Tasks t
    inner join Groups g on g.id = t.group_id
    inner join Statuses s on s.id = t.status_id
where 1=1
    and t.sys_status = 1
    and g.sys_status = 1
    and t.status_id = 4;

-- deleted
create view task_list_deleted as
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
    Tasks t
    inner join Groups g on g.id = t.group_id
    inner join Statuses s on s.id = t.status_id
where 1=1
    and (t.sys_status = 0 or g.sys_status = 0)

-- groups
create view group_list as 
select id, group_name from Groups where sys_status = 1;

create view group_list_deleted as 
select id, group_name from Groups where sys_status = 0;

-- status
create view status_list as 
select id, status_name from Statuses where sys_status = 1;

-- counts
create view task_counts as 
select 
    count(*) as count_all
    , sum(case when status_id = 1 then 1 else 0 end) as count_new
    , sum(case when status_id = 2 then 1 else 0 end) as count_in_progress
    , sum(case when status_id = 3 then 1 else 0 end) as count_on_hold
    , sum(case when status_id = 4 then 1 else 0 end) as count_completed
    , sum(case when status_id in (1, 2, 3) then 1 else 0 end) as count_ongoing
    , sum(case when date_due < current_date then 1 else 0 end) as count_overdue
from 
    Tasks t
    inner join Groups g on g.id = t.group_id
where 1=1
    and t.sys_status = 1
    and g.sys_status = 1;

create view group_counts as 
select 
    g.id
    , g.group_name
    , count(*) as count_all
    , sum(case when t.status_id = 1 then 1 else 0 end) as count_new
    , sum(case when t.status_id = 2 then 1 else 0 end) as count_in_progress
    , sum(case when t.status_id = 3 then 1 else 0 end) as count_on_hold
    , sum(case when t.status_id = 4 then 1 else 0 end) as count_completed
    , sum(case when t.status_id in (1, 2, 3) then 1 else 0 end) as count_ongoing
    , sum(case when t.date_due < current_date then 1 else 0 end) as count_overdue
from 
    Groups g
    left outer join Tasks t on t.group_id = g.id
        and t.sys_status = 1
where g.sys_status = 1
group by g.id, g.group_name;

create view status_counts as 
select 
    s.id
    , s.status_name
    , count(*) as count_all
    , sum(case when t.status_id = 1 then 1 else 0 end) as count_new
    , sum(case when t.status_id = 2 then 1 else 0 end) as count_in_progress
    , sum(case when t.status_id = 3 then 1 else 0 end) as count_on_hold
    , sum(case when t.status_id = 4 then 1 else 0 end) as count_completed
    , sum(case when t.status_id in (1, 2, 3) then 1 else 0 end) as count_ongoing
    , sum(case when t.date_due < current_date then 1 else 0 end) as count_overdue
from 
    Statuses s
    left outer join Tasks t on t.status_id = g.id
        and t.sys_status = 1
    inner join Groups g on g.id = t.group_id
        and g.sys_status = 1
where s.sys_status = 1
group by s.id, s.status_name
