-- public.movie_directors definition



CREATE TABLE public.movie_directors (
	movie_id int4 NOT NULL,
	director_id int4 NOT NULL,
	CONSTRAINT movie_directors_pkey PRIMARY KEY (movie_id, director_id),
	CONSTRAINT movie_directors_director_id_fkey FOREIGN KEY (director_id) REFERENCES public.directors(id) ON DELETE CASCADE,
	CONSTRAINT movie_directors_movie_id_fkey FOREIGN KEY (movie_id) REFERENCES public.movies(id) ON DELETE CASCADE
);