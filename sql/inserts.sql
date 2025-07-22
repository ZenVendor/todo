-- Initial inserts
insert or ignore into TaskStatus (id, status_name) values (1, 'New');
insert or ignore into TaskStatus (id, status_name) values (2, 'In Progress');
insert or ignore into TaskStatus (id, status_name) values (3, 'On Hold');
insert or ignore into TaskStatus (id, status_name) values (4, 'Completed');

insert or ignore into TaskGroup (id, group_name) values (1, 'General');
