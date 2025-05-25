-- public.tiers definition



CREATE TABLE public.tiers (
	id serial4 NOT NULL,
	"name" varchar(50) NOT NULL,
	min_points int4 NOT NULL,
	CONSTRAINT tiers_min_points_check CHECK ((min_points >= 0)),
	CONSTRAINT tiers_min_points_key UNIQUE (min_points),
	CONSTRAINT tiers_pkey PRIMARY KEY (id)
);