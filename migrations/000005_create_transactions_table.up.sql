-- public.transactions definition



CREATE TABLE public.transactions (
	id serial4 NOT NULL,
	user_id int4 NOT NULL,
	payment_id int4 NOT NULL,
	status_paid bool DEFAULT false NULL,
	total int4 NULL,
	transaction_date timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	invoice_number varchar(50) NULL,
	CONSTRAINT transactions_invoice_number_key UNIQUE (invoice_number),
	CONSTRAINT transactions_pkey PRIMARY KEY (id),
	CONSTRAINT transactions_total_check CHECK ((total >= 0)),
	CONSTRAINT transactions_payment_id_fkey FOREIGN KEY (payment_id) REFERENCES public.payments(id) ON DELETE CASCADE,
	CONSTRAINT transactions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE
);