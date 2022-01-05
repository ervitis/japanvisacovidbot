CREATE DATABASE japancovid ENCODING 'UTF8';

CREATE TABLE IF NOT EXISTS embassydates (
    date DATE NOT NULL,
    embassy VARCHAR(3) NOT NULL,
    PRIMARY KEY (date, embassy)
);

CREATE TABLE IF NOT EXISTS coviddata (
    id BIGSERIAL,
    datecovid DATE NOT NULL,
    date CHAR(8) NOT NULL,
    pcr BIGINT NOT NULL,
    positive BIGINT NOT NULL,
    symptom BIGINT NOT NULL,
    symptomless BIGINT NOT NULL,
    symtomConfirming BIGINT NOT NULL,
    hospitalize BIGINT NOT NULL,
    mild BIGINT NOT NULL,
    severe BIGINT NOT NULL,
    confirming BIGINT NOT NULL,
    waiting BIGINT NOT NULL,
    discharge BIGINT NOT NULL,
    death BIGINT NOT NULL,
    PRIMARY KEY (id, date)
)