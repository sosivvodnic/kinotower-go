-- Schema: kinotower (init)
-- This migration creates all tables required by the task.

BEGIN;

CREATE TABLE IF NOT EXISTS genders (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS categories (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(150) NOT NULL,
    parent_id  INT NULL,
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT fk_categories_parent
        FOREIGN KEY (parent_id) REFERENCES categories(id)
        ON UPDATE CASCADE
        ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_categories_parent_id ON categories(parent_id);
CREATE INDEX IF NOT EXISTS idx_categories_deleted_at ON categories(deleted_at);

CREATE TABLE IF NOT EXISTS countries (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id         SERIAL PRIMARY KEY,
    fio        VARCHAR(150) NOT NULL,
    birthday   DATE NULL,
    gender_id  INT NOT NULL,
    email      VARCHAR(50) NOT NULL,
    password   VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT uq_users_email UNIQUE (email),
    CONSTRAINT fk_users_gender
        FOREIGN KEY (gender_id) REFERENCES genders(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_users_gender_id ON users(gender_id);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

CREATE TABLE IF NOT EXISTS films (
    id             SERIAL PRIMARY KEY,
    name           VARCHAR(150) NOT NULL,
    country_id     INT NOT NULL,
    duration       INT NOT NULL,
    year_of_issue  INT NOT NULL,
    age            INT NOT NULL,
    link_img       VARCHAR(255) NULL,
    link_kinopoisk VARCHAR(255) NULL,
    link_video     VARCHAR(255) NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL,
    deleted_at     TIMESTAMPTZ NULL,
    CONSTRAINT fk_films_country
        FOREIGN KEY (country_id) REFERENCES countries(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    CONSTRAINT chk_films_duration_unsigned CHECK (duration >= 0),
    CONSTRAINT chk_films_age_unsigned CHECK (age >= 0),
    CONSTRAINT chk_films_year_of_issue CHECK (year_of_issue >= 1888)
);

CREATE INDEX IF NOT EXISTS idx_films_country_id ON films(country_id);
CREATE INDEX IF NOT EXISTS idx_films_deleted_at ON films(deleted_at);

CREATE TABLE IF NOT EXISTS categories_films (
    id          SERIAL PRIMARY KEY,
    category_id INT NOT NULL,
    film_id     INT NOT NULL,
    CONSTRAINT uq_categories_films_pair UNIQUE (category_id, film_id),
    CONSTRAINT fk_categories_films_category
        FOREIGN KEY (category_id) REFERENCES categories(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    CONSTRAINT fk_categories_films_film
        FOREIGN KEY (film_id) REFERENCES films(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_categories_films_category_id ON categories_films(category_id);
CREATE INDEX IF NOT EXISTS idx_categories_films_film_id ON categories_films(film_id);

CREATE TABLE IF NOT EXISTS reviews (
    id          SERIAL PRIMARY KEY,
    film_id     INT NOT NULL,
    user_id     INT NOT NULL,
    message     TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL,
    is_approved BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_at  TIMESTAMPTZ NULL,
    CONSTRAINT fk_reviews_film
        FOREIGN KEY (film_id) REFERENCES films(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_reviews_user
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_reviews_film_id ON reviews(film_id);
CREATE INDEX IF NOT EXISTS idx_reviews_user_id ON reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_reviews_deleted_at ON reviews(deleted_at);

CREATE TABLE IF NOT EXISTS ratings (
    id         SERIAL PRIMARY KEY,
    film_id    INT NOT NULL,
    user_id    INT NOT NULL,
    ball       INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT uq_ratings_pair UNIQUE (film_id, user_id),
    CONSTRAINT fk_ratings_film
        FOREIGN KEY (film_id) REFERENCES films(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_ratings_user
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    CONSTRAINT chk_ratings_ball_unsigned CHECK (ball >= 0),
    CONSTRAINT chk_ratings_ball_range CHECK (ball >= 1 AND ball <= 5)
);

CREATE INDEX IF NOT EXISTS idx_ratings_film_id ON ratings(film_id);
CREATE INDEX IF NOT EXISTS idx_ratings_user_id ON ratings(user_id);

COMMIT;
