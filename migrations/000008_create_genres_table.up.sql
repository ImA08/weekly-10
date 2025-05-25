-- public.genres definition



CREATE TABLE public.genres (
	id serial4 NOT NULL,
	genre varchar(50) NOT NULL,
	CONSTRAINT genres_pkey PRIMARY KEY (id)
);