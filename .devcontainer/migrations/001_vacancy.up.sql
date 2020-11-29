CREATE SCHEMA vacancy;

CREATE TYPE vacancy.STATUS AS enum (
  'none',
  'draft',
  'active',
  'inactive'
);

CREATE TABLE vacancy.vacancy (
  id            TEXT,
  template_id   TEXT            NOT NULL,
  title         TEXT            NOT NULL,
  status        vacancy.STATUS  NOT NULL,
  area          TEXT,
  department    TEXT,
  duties        TEXT[],
  requirements  TEXT[],
  experience    int,
  created       TIMESTAMP       NOT NULL,
  updated       TIMESTAMP       NOT NULL,
  
  CONSTRAINT pk_vacancy__id PRIMARY KEY (id)
);

CREATE INDEX ix_vacancy__templated_id ON vacancy.vacancy (template_id);
CREATE INDEX ix_vacancy__title        ON vacancy.vacancy (title);
CREATE INDEX ix_vacancy__status       ON vacancy.vacancy (status);
CREATE INDEX ix_vacancy__area     ON vacancy.vacancy (area);
CREATE INDEX ix_vacancy__ldepartment  ON vacancy.vacancy (department);
CREATE INDEX ix_vacancy__created      ON vacancy.vacancy (created);

CREATE TABLE vacancy.skill (
  vacancy_id  TEXT,
  title       TEXT NOT NULL,
  important   BOOLEAN NOT NULL,
  
  CONSTRAINT pk_skill__id PRIMARY KEY (vacancy_id, title),
  CONSTRAINT fk_skill__vacancy_id FOREIGN KEY (vacancy_id) REFERENCES vacancy.vacancy (id)
);

CREATE INDEX ix_skill__title     ON vacancy.skill (title);
CREATE INDEX ix_skill__important ON vacancy.skill (important);
