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
  channelCreatorUserID int not null,
  foreign key(channelCreatorUserID) references users(id),
  channelEditedAt datetime,
  unique(channelName)
);

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
  foreign key(mCreatorUserID) references users(id),
  mEditedAt datetime not null
);

create table if not exists message_reaction (
  mrID int primary key auto_increment not null,
  mrMessageID int not null,
  mrUserID int not null,
  mrReactionCode varchar(255) not null,
  foreign key(mrMessageID) references messages(mMessageID),
  foreign key(mrUserID) references users(id),
  unique key(mrMessageID, mrUserID, mrReactionCode)
);

create table if not exists star_message (
  smUserID int not null,
  smMessageID int not null,
  foreign key(smUserID) references users(id),
  foreign key(smMessageID) references messages(mMessageID),
  primary key (smUserID, smMessageID)
);

insert into users (email, passHash, username, firstName, lastName, photoURL)
values ("admin@system.com", "\0", "system", "", "", "");

insert into channels (channelName, channelDescription, channelPrivate, channelCreatedAt, channelCreatorUserID, channelEditedAt)
values ('general', 'general channel', false, null, 1, null);

insert into channel_user (cuUserID, cuChannelID)
values (1, 1);