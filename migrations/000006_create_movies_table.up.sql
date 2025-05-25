-- public.movies definition



CREATE TABLE public.movies (
	id serial4 NOT NULL,
	title varchar(255) NOT NULL,
	synopsis text NULL,
	duration int4 NOT NULL,
	release_date date NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT movies_pkey PRIMARY KEY (id)
);