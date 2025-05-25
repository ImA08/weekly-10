-- public.directors definition

-- Drop table

-- DROP TABLE public.directors;

CREATE TABLE public.directors (
	id serial4 NOT NULL,
	"name" varchar(255) NOT NULL,
	CONSTRAINT directors_pkey PRIMARY KEY (id)
);