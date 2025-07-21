begin transaction;
with old_tasks as (
select 
    id
    , description as summary
    , priority 
    , due as date_due
    , completed as date_completed
    , null as description
    , null as closing_comment
    , case when done = 1 then 4 else 2 end as status_id
    , group_id
    , null as parent_id
    , created as sys_created
    , updated as sys_updated
    , 1 as sys_status
from tasklist_old
)
insert into Task 
select * from old_tasks;
commit transaction;

begin transaction;
with old_groups as (
select 
    id
    , name as group_name
    , created as sys_created
    , updated as sys_updated
    , 1 as sys_status
from taskgroup_old
where id <> 1
)
insert into TaskGroup 
select * from old_groups;
commit transaction;
