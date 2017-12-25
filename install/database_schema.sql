/*
    Setting up the database is covered in the README.md file
*/

USE isaac;

/*
    We have to disable foreign key checks so that we can drop the tables;
    this will only disable it for the current session
*/
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS tournament_races;
CREATE TABLE tournament_races (
    id                    INT            NOT NULL  PRIMARY KEY  AUTO_INCREMENT,
    /* PRIMARY KEY automatically creates a UNIQUE constraint */
    tournament_name       NVARCHAR(100)  NOT NULL,
    racer1                INT            NOT NULL,
    racer2                INT            NOT NULL,
    channel_id            NVARCHAR(100)  NOT NULL,
    challonge_id          NVARCHAR(100)  NOT NULL,
    bracket_round         NVARCHAR(10)   NOT NULL,
    state                 NVARCHAR(50)   NOT NULL, /* definitions are listed in the "Race" struct */
    datetime_created      TIMESTAMP      NOT NULL  DEFAULT NOW(),
    datetime_scheduled    TIMESTAMP      NULL      DEFAULT NULL,
    caster                INT            NULL      DEFAULT NULL,
    caster_p1             INT            NOT NULL  DEFAULT 0,
    caster_p2             INT            NOT NULL  DEFAULT 0,
    active_player         INT            NOT NULL  DEFAULT 1,
    characters_remaining  NVARCHAR(500)  NOT NULL,
    characters            NVARCHAR(500)  NOT NULL  DEFAULT "",
    builds_remaining      NVARCHAR(500)  NOT NULL,
    builds                NVARCHAR(500)  NOT NULL  DEFAULT "",
    racer1_bans           INT            NOT NULL,
    racer2_bans           INT            NOT NULL,
    racer1_vetos          INT            NOT NULL,
    racer2_vetos          INT            NOT NULL,
    num_voted             INT            NOT NULL  DEFAULT 0,
    score                 NVARCHAR(10)   NOT NULL  DEFAULT "0-0",
    FOREIGN KEY (racer1) REFERENCES tournament_racers (id) ON DELETE CASCADE,
    FOREIGN KEY (racer2) REFERENCES tournament_racers (id) ON DELETE CASCADE,
    FOREIGN KEY (caster) REFERENCES tournament_racers (id) ON DELETE CASCADE
);
CREATE INDEX tournament_races_index_channel_id ON tournament_races (channel_id);

DROP TABLE IF EXISTS tournament_racers;
CREATE TABLE tournament_racers (
    id          INT            NOT NULL  PRIMARY KEY  AUTO_INCREMENT,
    /* PRIMARY KEY automatically creates a UNIQUE constraint */
    discord_id  NVARCHAR(100)  NOT NULL  UNIQUE,
    username    NVARCHAR(100)  NOT NULL,
    timezone    NVARCHAR(100)  NULL      DEFAULT NULL,
    /* the TZ column of: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones */
    stream_url  NVARCHAR(255)  NULL      DEFAULT NULL
);
CREATE INDEX tournament_racers_index_discord_id ON tournament_racers (discord_id);
CREATE INDEX tournament_racers_index_username ON tournament_racers (username);
