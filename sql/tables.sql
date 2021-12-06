CREATE DATABASE japancovid ENCODING 'UTF8';

CREATE TABLE IF NOT EXISTS embassydates (
    date DATE NOT NULL,
    embassy VARCHAR(3) NOT NULL,
    PRIMARY KEY (date, embassy)
);