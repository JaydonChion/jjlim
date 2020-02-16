drop table sessions;
drop table users;


create table users (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  name       varchar(255),
  mnum       varchar(255),
  email      varchar(255) not null unique,
  password   varchar(255),
  created_at timestamp not null,
  phone int default 0,
  dob timestamp not null, 
  gender  varchar(64) not null,
  verification_status bool default 'false',
  cupSaved int default 0
);

create table sessions (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  user_id    integer references users(id),
  created_at timestamp not null
);

create table rentalData(

  incident  serial primary key,
  user_id    integer references users(id),
  created_at timestamp not null,
  item_id int,
  location_id int

);

create table donorData(

      location varchar(255) not null,
      name varchar(255) not null,
      user_id integer references users(id),
      phoneNum int,
      Mnum varchar(255),
      halal varchar(255),
      amountNum int,
      descrp varchar(255),
      created_at timestamp not null

);

create table predictData(

  tpeople int, 
  nmale int, 
  nfemale int, 
  age int,
  timeNum int,
  foodleft int
);
