-- schema for Tasks database
create table tasks (
    id varchar(25) primary key not null,
    title varchar(255) not null,
    completed bool not null default 0,
    createdAt datetime not null,
    modifiedAt datetime null
);

create table tags (
    taskID varchar(25),
    tag varchar(64) not null,
    foreign key (taskID) references tasks(id)
)
