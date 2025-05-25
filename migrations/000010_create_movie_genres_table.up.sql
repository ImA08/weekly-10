-- public.movie_genres definition



CREATE TABLE public.movie_genres (
	movie_id int4 NOT NULL,
	genre_id int4 NOT NULL,
	CONSTRAINT movie_genres_pkey PRIMARY KEY (movie_id, genre_id),
	CONSTRAINT movie_genres_genre_id_fkey FOREIGN KEY (genre_id) REFERENCES public.genres(id) ON DELETE CASCADE,
	CONSTRAINT movie_genres_movie_id_fkey FOREIGN KEY (movie_id) REFERENCES public.movies(id) ON DELETE CASCADE
);