-- all
drop view if exists task_list_all;
create view if not exists task_list_all as
select 
    t.id 
    , t.summary
    , t.priority
    , t.date_due
    , t.date_completed
    , t.description
    , t.closing_comment
    , t.status
    , t.project_id
    , p.project_name
    , t.parent_id
    , t.sys_created
    , t.sys_updated
    , t.sys_status
from
    Task t
    inner join Project p on p.id = t.project_id
where 1=1
    and t.sys_status = 1;

-- new
drop view if exists task_list_new;
create view if not exists task_list_new as
select 
    t.id 
    , t.summary
    , t.priority
    , t.date_due
    , t.date_completed
    , t.description
    , t.closing_comment
    , t.status
    , t.project_id
    , p.project_name
    , t.parent_id
    , t.sys_created
    , t.sys_updated
    , t.sys_status
from
    Task t
    inner join Project p on p.id = t.project_id
where 1=1
    and t.sys_status = 1
    and t.status = 1;

-- open
drop view if exists task_list_open;
create view if not exists task_list_open as
select 
    t.id 
    , t.summary
    , t.priority
    , t.date_due
    , t.date_completed
    , t.description
    , t.closing_comment
    , t.status
    , t.project_id
    , p.project_name
    , t.parent_id
    , t.sys_created
    , t.sys_updated
    , t.sys_status
from
    Task t
    inner join Project p on p.id = t.project_id
where 1=1
    and t.sys_status = 1
    and t.status in (1, 2, 3);

-- in_progress
drop view if exists task_list_in_progress;
create view if not exists task_list_in_progress as
select 
    t.id 
    , t.summary
    , t.priority
    , t.date_due
    , t.date_completed
    , t.description
    , t.closing_comment
    , t.status
    , t.project_id
    , p.project_name
    , t.parent_id
    , t.sys_created
    , t.sys_updated
    , t.sys_status
from
    Task t
    inner join Project p on p.id = t.project_id
where 1=1
    and t.sys_status = 1
    and t.status = 2;

-- on hold
drop view if exists task_list_on_hold;
create view if not exists task_list_on_hold as
select 
    t.id 
    , t.summary
    , t.priority
    , t.date_due
    , t.date_completed
    , t.description
    , t.closing_comment
    , t.status
    , t.project_id
    , p.project_name
    , t.parent_id
    , t.sys_created
    , t.sys_updated
    , t.sys_status
from
    Task t
    inner join Project p on p.id = t.project_id
where 1=1
    and t.sys_status = 1
    and t.status = 3;

-- completed
drop view if exists task_list_completed;
create view if not exists task_list_completed as
select 
    t.id 
    , t.summary
    , t.priority
    , t.date_due
    , t.date_completed
    , t.description
    , t.closing_comment
    , t.status
    , t.project_id
    , p.project_name
    , t.parent_id
    , t.sys_created
    , t.sys_updated
    , t.sys_status
from
    Task t
    inner join Project p on p.id = t.project_id
where 1=1
    and t.sys_status = 1
    and t.status = 4;

-- deleted
drop view if exists task_list_deleted;
create view if not exists task_list_deleted as
select 
    t.id 
    , t.summary
    , t.priority
    , t.date_due
    , t.date_completed
    , t.description
    , t.closing_comment
    , t.status
    , t.project_id
    , p.project_name
    , t.parent_id
    , t.sys_created
    , t.sys_updated
    , t.sys_status
from
    Task t
    inner join Project p on p.id = t.project_id
where 1=1
    and t.sys_status = 0;

-- overdue
drop view if exists task_list_overdue;
create view if not exists task_list_overdue as
select 
    t.id 
    , t.summary
    , t.priority
    , t.date_due
    , t.date_completed
    , t.description
    , t.closing_comment
    , t.status
    , t.project_id
    , p.project_name
    , t.parent_id
    , t.sys_created
    , t.sys_updated
    , t.sys_status
from
    Task t
    inner join Project p on p.id = t.project_id
where 1=1
    and t.sys_status = 1
    and t.status in (1, 2, 3)
    and t.date_due is not null
    and t.date_due < current_date;


-- groups
drop view if exists project_list ;
create view if not exists project_list as 
select id, project_name from Project where sys_status = 1;

-- counts
drop view if exists status_counts ;
create view if not exists status_counts as 
select 
    count(*) as count_all
    , sum(case when status = 1 then 1 else 0 end) as count_new
    , sum(case when status = 2 then 1 else 0 end) as count_in_progress
    , sum(case when status = 3 then 1 else 0 end) as count_on_hold
    , sum(case when status = 4 then 1 else 0 end) as count_completed
    , sum(case when status in (1, 2, 3) then 1 else 0 end) as count_open
    , sum(case 
        when status in (1, 2, 3) 
            and date_due is not null
            and date_due < current_date
        then 1 
        else 0 
    end) as count_overdue
from Task
where sys_status = 1;

drop view if exists project_counts ;
create view if not exists project_counts as 
select 
    project_id
    , count(*) as count_all
    , sum(case when status = 1 then 1 else 0 end) as count_new
    , sum(case when status = 2 then 1 else 0 end) as count_in_progress
    , sum(case when status = 3 then 1 else 0 end) as count_on_hold
    , sum(case when status = 4 then 1 else 0 end) as count_completed
    , sum(case when status in (1, 2, 3) then 1 else 0 end) as count_open
    , sum(case 
        when status in (1, 2, 3) 
            and date_due is not null
            and date_due < current_date
        then 1 
        else 0 
    end) as count_overdue
from Task
where sys_status = 1
group by project_id
order by project_id;
