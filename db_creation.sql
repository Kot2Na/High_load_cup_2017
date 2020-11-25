create database if not exists travel_service;
create user if not exists kot@localhost identified by 'pass';
grant all privileges on travel_service.* to kot@localhost;

use travel_service;

create table users (
	user_id int unsigned not null,
    email varchar(100) not null,
    first_name varchar(50) not null,
    last_name varchar(50) not null,
    gender set('m', 'f') not null,
    birth_date int not null,
    
    primary key(user_id)
);

create table locations (
	loc_id int unsigned not null,
    place text not null,
    country varchar(50) not null,
    city varchar(50) not null,
    distance int unsigned not null,
    
    primary key (loc_id)
);

create table visits (
	visit_id int unsigned not null,
    user_id int unsigned not null,
    loc_id int unsigned not null,
    visit_date int not null,
    mark tinyint unsigned not null,
    
    foreign key (user_id) references users (user_id) on delete cascade, 
    foreign key (loc_id) references locations (loc_id) on delete cascade,
    primary key (visit_id)
);