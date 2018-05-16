create table if not exists users (
  id int primary key auto_increment not null,
  email varchar(255) not null,
  passHash binary(60) not null,
  username varchar(255) not null,
  firstName varchar(35) not null,
  lastName varchar(35) not null,
  photoURL varchar(2083) not null,
  unique(email),
  unique(username)
);

create table if not exists channels (
  id int primary key auto_increment not null,
  name varchar(255) not null,
  description varchar(1024),
  private bool not null default false,
  createdAt datetime,
  creatorUserID int,
  editedAt datetime,
  unique(name)
);

insert into channels (name, description, private, createdAt, creatorUserID, editedAt)
values ('general', 'general channel', false, null, null, null);

create table if not exists channel_user (
  id int primary key auto_increment not null,
  userID int not null,
  channelID int not null,
  foreign key(userID) references users(id),
  foreign key(channelID) references channels(id)
);

create table if not exists messages (
  id int primary key auto_increment not null,
  channelID int not null,
  foreign key(channelID) references channels(id),
  body varchar(4000) not null,
  createdAt datetime not null,
  creatorUserID int not null,
  editedAt datetime not null
);