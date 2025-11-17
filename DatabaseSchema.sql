-- Formally verified E-Voting using Dafny
-- Copyright (C) 2025 Authors Gruppe EinS
-- This program is free software: you can redistribute it and/or modify
-- it under the terms of the GNU Affero General Public License as
-- published by the Free Software Foundation, either version 3 of the
-- License, or (at your option) any later version.
-- This program is distributed in the hope that it will be useful,
-- but WITHOUT ANY WARRANTY; without even the implied warranty of
-- MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
-- GNU Affero General Public License for more details.
-- You should have received a copy of the GNU Affero General Public License
-- along with this program.  If not, see <https://www.gnu.org/licenses/>.


# V1

create table wahlleiter (
	id int NOT NULL primary key,
	username varchar(200) unique,
	email varchar(200),
	password varchar(64),
	salt varchar(24)
);

CREATE SEQUENCE seq_wahlleiter;

create table election (
	id int not null primary key,
	name varchar(200),
	beschreibung TEXT,
	wahlleiter_id int,
	created_at DATETIME,
	end_time DATETIME,
	type int, # Werte Enum 1 - 4, erweiterbar
	is_active bool,
	constraint fk_wahlleiter_id foreign key (wahlleiter_id) references wahlleiter(id)
);

CREATE SEQUENCE seq_election;

create table candidates (
	id int not null primary key,
	election_id int,
	name varchar(255),
	constraint fk_candidates foreign key (election_id) references election(id)
);

create SEQUENCE seq_candidates;

create table vote_header (
	id int not null primary key,
	election_id int,
	constraint fk_voteheader foreign key (election_id) references election(id)
);

create sequence seq_vote_header;

create table single_vote_approval (
	vote_id int, 
	candidate_id int, 
	approved bool,
	primary key(vote_id, candidate_id),
	constraint fk_single_vote_approval_vote foreign key (vote_id) references vote_header(id),
	constraint fk_single_vote_approval_candidate foreign key (candidate_id) references candidates(id)
);

# V2
create sequence seq_wahltoken;

create table wahltoken (
	id int not null primary key,
	election_id int,
	token varchar(32),
	voted bool, # irrelevant für offene Wahl
	constraint fk_wahltoken_election foreign key (election_id) references election(id)
);

alter table election add password varchar(64);
alter table election add salt varchar(24);
alter table election add hmac varchar(64);
alter table election add open_wahl bool;

# V3 - unified Voting API
create sequence seq_Vote_Type2;

# Für Rejection, Scored und IRV
create table single_vote_type2 (
	vote_id int, 
	candidate_id int, 
	info int,
	primary key(vote_id, candidate_id),
	constraint fk_single_vote_type2_vote foreign key (vote_id) references vote_header(id),
	constraint fk_single_vote_type2_candidate foreign key (candidate_id) references candidates(id)
);
