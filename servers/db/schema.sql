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
  channelID int primary key auto_increment not null,
  channelName varchar(255) not null,
  channelDescription varchar(1024),
  channelPrivate bool not null default false,
  channelCreatedAt datetime,
  channelCreatorUserID int,
  channelEditedAt datetime,
  unique(channelName)
);

insert into channels (channelName, channelDescription, channelPrivate, channelCreatedAt, channelCreatorUserID, channelEditedAt)
values ('general', 'general channel', false, null, null, null);

create table if not exists channel_user (
  cuUserID int not null,
  cuChannelID int not null,
  foreign key(cuUserID) references users(id),
  foreign key(cuChannelID) references channels(channelID),
  primary key (cuUserID, cuChannelID)
);

create table if not exists messages (
  mMessageID int primary key auto_increment not null,
  mChannelID int not null,
  foreign key(mChannelID) references channels(channelID),
  mBody varchar(4000) not null,
  mCreatedAt datetime not null,
  mCreatorUserID int not null,
  mEditedAt datetime not null
);