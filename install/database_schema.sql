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
    tournament_name       NVARCHAR(500)  NOT NULL,
    racer1                INT            NOT NULL, /* The "tournament_users" database ID */
    racer1_challonge_id   INT            NOT NULL, /* The "participant" ID; needed to automatically set the winner through the Challonge API */
    racer2                INT            NOT NULL, /* The "tournament_users" database ID */
    racer2_challonge_id   INT            NOT NULL, /* The "participant" ID; needed to automatically set the winner through the Challonge API */
    channel_id            NVARCHAR(100)  NOT NULL, /* The Discord channel ID that was automatically created for this race */
    channel_name          NVARCHAR(500)  NOT NULL,
    challonge_url         NVARCHAR(100)  NOT NULL, /* The suffix of the Challonge URL for this tournament */
    challonge_match_id    NVARCHAR(100)  NOT NULL,
    bracket_round         NVARCHAR(10)   NOT NULL,
    state                 NVARCHAR(50)   NOT NULL, /* Definitions are listed in the "Race" struct */
    datetime_created      TIMESTAMP      NOT NULL  DEFAULT NOW(),
    datetime_scheduled    TIMESTAMP      NULL      DEFAULT NULL,
    first_picker          INT            NOT NULL  DEFAULT 1,
    active_racer          INT            NOT NULL  DEFAULT 1,
    characters_remaining  NVARCHAR(500)  NOT NULL,
    characters            NVARCHAR(500)  NOT NULL  DEFAULT "",
    builds_remaining      NVARCHAR(500)  NOT NULL,
    builds                NVARCHAR(500)  NOT NULL  DEFAULT "",
    racer1_bans           INT            NOT NULL,
    racer2_bans           INT            NOT NULL,
    racer1_vetos          INT            NOT NULL,
    racer2_vetos          INT            NOT NULL,
    num_voted             INT            NOT NULL  DEFAULT 0,
    score                 NVARCHAR(10)   NULL      DEFAULT NULL, /* e.g. "3-2" */
    FOREIGN KEY (racer1) REFERENCES tournament_users (id) ON DELETE CASCADE,
    FOREIGN KEY (racer2) REFERENCES tournament_users (id) ON DELETE CASCADE
);
CREATE INDEX tournament_races_index_channel_id ON tournament_races (channel_id);

DROP TABLE IF EXISTS tournament_users;
CREATE TABLE tournament_users (
    id                       INT            NOT NULL  PRIMARY KEY  AUTO_INCREMENT,
    /* PRIMARY KEY automatically creates a UNIQUE constraint */
    discord_id               NVARCHAR(100)  NOT NULL  UNIQUE,
    username                 NVARCHAR(100)  NOT NULL,
    timezone                 NVARCHAR(100)  NULL      DEFAULT NULL,
    /* The TZ column of: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones */
    stream_url               NVARCHAR(255)  NULL      DEFAULT NULL,
    caster_always_ok         TINYINT(1)     NOT NULL  DEFAULT 0
);
CREATE INDEX tournament_users_index_discord_id ON tournament_users (discord_id);
CREATE INDEX tournament_users_index_username ON tournament_users (username);

DROP TABLE IF EXISTS tournament_casts;
CREATE TABLE tournament_casts (
    id             INT           NOT NULL  PRIMARY KEY  AUTO_INCREMENT,
    race_id        INT           NOT NULL, /* The "tournament_races" database ID */
    caster         INT           NOT NULL, /* The "tournament_users" database ID */
    r1_permission  INT           NOT NULL  DEFAULT 0, /* Whether or not racer 1 has given permission to this caster */
    r2_permission  INT           NOT NULL  DEFAULT 0, /* Whether or not racer 2 has given permission to this caster */
    language       NVARCHAR(50)  NOT NULL, /* Equal to "en", "fr", etc. */
    FOREIGN KEY (race_id) REFERENCES tournament_races (id) ON DELETE CASCADE,
    FOREIGN KEY (caster) REFERENCES tournament_users (id) ON DELETE CASCADE,
    UNIQUE(race_id, caster), /* The same person cannot cast the same race more than once */
    UNIQUE(race_id, language) /* There cannot be two casts of the same race in the same language */
);
