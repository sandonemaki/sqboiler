CREATE TABLE user_favorite_movies (
    user_id INTEGER REFERENCES users(id),
    movie_id INTEGER REFERENCES movies(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, movie_id)
);
