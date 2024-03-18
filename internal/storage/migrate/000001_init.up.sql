CREATE TABLE actors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    gender VARCHAR(10) NOT NULL,
    birth_date DATE NOT NULL
);

CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(150) NOT NULL,
    description TEXT NOT NULL,
    release_date DATE NOT NULL,
    rating NUMERIC(2,1) NOT NULL CHECK (rating >= 0 AND rating <= 10),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE movie_actors (
    movie_id INTEGER NOT NULL REFERENCES movies(id),
    actor_id INTEGER NOT NULL REFERENCES actors(id),
    PRIMARY KEY (movie_id, actor_id)
);

CREATE INDEX idx_movies_title ON movies (title);
CREATE INDEX idx_movies_release_date ON movies (release_date);
CREATE INDEX idx_movies_rating ON movies (rating);
CREATE INDEX idx_actors_name ON actors (name);