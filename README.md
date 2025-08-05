# TODO CLI app
This is a command line TODO program.
It does one thing at a time, by design. 

## TODO
* Add project functions (counts, rename, delete)
* Create install script
* There's more but it escapes my memory...

## Config
* Default config file: todo_config.yml
* Default db file: todo.db

The program looks for the config file in current dir, $XDG_CONFIG_HOME/todo or $HOME/.config/todo, in that order.

### Config parameters:
* dbLocation: $HOME/.config/todo/
* dbName: todo.db
* dateFormat: 2006-01-02            - this uses go date format.
* defaultProject: General           - name of the default project
* projectNameLength: 24             - max length of project names
* summaryLength: 255                - max length of summary
* editor: /usr/bin/vim              - text editor to use for description and closing comment

## Usage

todo [verb] [required_value] [args] [kwargs]
	
Executed without arguments defaults to "list --open"
Due date and parent ID can be removd by providing empty value (--due= or --parent= )
Values for key-value arguments are provided after =, e.g. --project=Todo
--description and --comment dod not take a value but open text editor, set in config file.

```
add <summary>               Adds new task
    --description, --desc
    --due                   [=date]
    --project, --proj       [=project_name]
    --priority, --pty       [=number|keyword]
    --parent, --pid         [=task_id]

complete <task_id>          Sets task completed, including child tasks
    --comment

count                       Displays count for selected status, useful 
                            for prompt, except --all which is default 
                            and displays list of statuses
    --all
    --completed
    --deleted
    --inprog
    --new
    --onhold
    --open
    --overdue

delete <task_id>            Soft deletes task and unlinks child tasks,
                            or also deletes them with --all
    --all

help                        Displays help

hold                        Puts task on hold
    --description, --desc

list                        Lists tasks, defaults to --open            
    --completed
    --deleted
    --inprog
    --new
    --onhold
    --open
    --overdue

reopen <task_id>            For completed tasks - switches to New 
                            status and appends closing comment to
                            description
    --description

show <task_id>              Displays task details
    
start <task_id>             Switches status to In progress
    --description

update <task_id>            Updates task
    --comment
    --description, --desc
    --due                   [=date]
    --project, --proj       [=project_name]
    --priority, --pty       [=number|keyword]
    --parent, --pid         [=task_id]
    --summary, --sum        [=summary]

start <task_id>             Undeletes task

version                     Displays version
```

## Examples
```
todo
todo add "New task" --parent=1
todo add "New task" --due=2024-08-13 --proj=Todo --priority=low
todo list --overdue
todo count --all
todo update 15 --summary="new summary" 
todo start 5 --description
todo update 10 --due= 
todo complete 12
todo reopen 3
todo delete 5
``` 
## Values
* task_id: integer
* summary: string with configurable maximum length (default 255)
* due: date in the format 2006-01-02 or 20060102 or anything in between (dashes are optional)
* project_name: string with configurable maximum length (default 24)  
* priority: number or keyword (none|low|medium|high|critical). Number ranges correspond to the keywords
    * critical: 1-9
    * high: 10-99
    * medium: 100-999
    * low: 1000-9999
    * none (reminder): 10000+

