create table users (
  id int primary key not null,
  email varchar(255) not null,
  passHash binary(60) not null,
  username varchar(255) not null,
  firstName varchar(35) not null,
  lastName varchar(35) not null,
  photoURL varchar(2083) not null,
);