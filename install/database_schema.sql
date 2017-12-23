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
    id                  INT            NOT NULL  PRIMARY KEY  AUTO_INCREMENT, /* PRIMARY KEY automatically creates a UNIQUE constraint */
    racer1              INT            NOT NULL,
    racer2              INT            NOT NULL,
    channel_id          NVARCHAR(100)  NOT NULL,
    bracket_round       NVARCHAR(10)   NOT NULL,
    state               INT            NOT NULL  DEFAULT 0, /* 0 is "not scheduled", 1 is "scheduled", etc. */
    datetime_created    TIMESTAMP      NOT NULL  DEFAULT NOW(),
    datetime_scheduled  TIMESTAMP      NULL      DEFAULT NULL,
    caster              INT            NULL      DEFAULT NULL,
    caster_p1           INT            NOT NULL  DEFAULT 0,
    caster_p2           INT            NOT NULL  DEFAULT 0,
    active_player       INT            NOT NULL  DEFAULT 1,
    characters          NVARCHAR(255)  NOT NULL,
    builds              NVARCHAR(500)  NOT NULL,
    racer1_veto         INT            NOT NULL  DEFAULT 0,
    racer2_veto         INT            NOT NULL  DEFAULT 0,
    FOREIGN KEY (racer1) REFERENCES tournament_racers (id) ON DELETE CASCADE,
    FOREIGN KEY (racer2) REFERENCES tournament_racers (id) ON DELETE CASCADE,
    FOREIGN KEY (caster) REFERENCES tournament_racers (id) ON DELETE CASCADE
);
CREATE INDEX tournament_races_index_channel_id ON tournament_races (channel_id);

DROP TABLE IF EXISTS tournament_racers;
CREATE TABLE tournament_racers (
    id          INT            NOT NULL  PRIMARY KEY  AUTO_INCREMENT, /* PRIMARY KEY automatically creates a UNIQUE constraint */
    discord_id  NVARCHAR(100)  NOT NULL  UNIQUE,
    username    NVARCHAR(100)  NOT NULL,
    timezone    NVARCHAR(100)  NULL      DEFAULT NULL, /* the number of hours adjusted from GMT */
    stream_url  NVARCHAR(255)  NULL      DEFAULT NULL
);
CREATE INDEX tournament_racers_index_discord_id ON tournament_racers (discord_id);
CREATE INDEX tournament_racers_index_username ON tournament_racers (username);
