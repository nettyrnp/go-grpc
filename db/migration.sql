
DROP TABLE IF EXISTS public.people;

CREATE TABLE public.people
(
  id INTEGER PRIMARY KEY NOT NULL,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(100) NOT NULL,
  mobile_number VARCHAR(30) NOT NULL
);
