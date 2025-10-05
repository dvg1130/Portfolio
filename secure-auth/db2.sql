
-- make sure var mathc in queries


CREATE TABLE IF NOT EXISTS users (
    uuid     CHAR(36) NOT NULL DEFAULT (UUID()),
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role     VARCHAR(255) NOT NULL DEFAULT 'basic'

);

CREATE TABLE IF NOT EXISTS users (
    uuid     CHAR(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
    ->     username VARCHAR(255) NOT NULL PRIMARY KEY,
    ->     password VARCHAR(255) NOT NULL,
    ->     role     VARCHAR(255) NOT NULL DEFAULT 'basic'
    -> );

    CREATE DATABASE IF NOT EXISTS records;
USE records;


CREATE TABLE snakes (
    suid        CHAR(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
    uuid  CHAR(36) NOT NULL,
    sid        VARCHAR(100) NOT NULL,
    species     VARCHAR(100) NOT NULL,
    sex         CHAR(1),
    age   CHAR(3),
    genes       CHAR(255),
    notes       CHAR(255)

);

CREATE TABLE feeding (
    suid        CHAR(36) NOT NULL PRIMARY KEY,
    uuid  CHAR(36) NOT NULL,
    sid        VARCHAR(100) NOT NULL,
    feed_date   DATE NOT NULL,
    prey_type   VARCHAR(100),
    prey_size    VARCHAR(50),
    notes       TEXT(255),
    FOREIGN KEY (suid) REFERENCES snakes(suid)
);

CREATE TABLE health (
    suid        CHAR(36) NOT NULL,
    uuid  CHAR(36) NOT NULL,
    sid         VARCHAR(100) NOT NULL,
    check_date  DATE NOT NULL,
    weight      VARCHAR(36),
    length      VARCHAR(36),
    topic     VARCHAR(50),
    notes       TEXT(255),
    PRIMARY KEY (suid, check_date),
    FOREIGN KEY (suid) REFERENCES snakes(suid)
);

-- breeding

CREATE TABLE breeding (
    breeding_uuid        CHAR(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
    uuid           CHAR(36) NOT NULL,
    event_id    CHAR(36) NOT NULL,
    female_suid    CHAR(36) NOT NULL,
    male1_suid      CHAR(36) NOT NULL,
    male2_suid      CHAR(36),
    male3_suid      CHAR(36),
    male4_suid      CHAR(36),
   female_sid    CHAR(36) NOT NULL,
    male1_sid      CHAR(36) NOT NULL,
    male2_sid      CHAR(36),
    male3_sid      CHAR(36),
    male4_sid      CHAR(36),
    breeding_year  DATE NOT NULL,
    breeding_season  CHAR(36),
    female_weight  CHAR(36),
    male1_weight      CHAR(36),
    male2_weight      CHAR(36),
    male3_weight      CHAR(36),
    male4_weight      CHAR(36),
    cooling_start  DATE,
    cooling_end  DATE,
    warming_start  DATE,
    warming_end  DATE,
    pairing1_date  DATE,
    pairing2_date  DATE,
    pairing3_date  DATE,
    pairing4_date  DATE,
    gravid_date  DATE,
    lay_date DATE,
    clutch_size INT,
    clutch_survive VARCHAR(36),
    outcome        VARCHAR(50),
    notes          TEXT(255),
    FOREIGN KEY (female_suid)      REFERENCES snakes(suid),
    FOREIGN KEY (male1_suid) REFERENCES snakes(suid)

);

-- offspring

CREATE TABLE offspring (
    suid        CHAR(36) NOT NULL,
    uuid  CHAR(36) NOT NULL,
    breeding_uuid   CHAR(36) NOT NULL,
    sid         VARCHAR(100) NOT NULL,
    hatch_date  DATE NOT NULL,
    mother_suid      CHAR(36) NOT NULL,
    mother_sid      CHAR(36) NOT NULL,
    father_suid      CHAR(36) NOT NULL,
    father_sid      CHAR(36) NOT NULL,
    baby_genes       CHAR(255),
    mother_genes       CHAR(255),
    father_genes       CHAR(255),
    notes       TEXT(255),
    PRIMARY KEY (suid, check_date),
    FOREIGN KEY (suid) REFERENCES snakes(suid)
);