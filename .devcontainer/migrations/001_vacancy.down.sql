DROP TABLE IF EXISTS public.skill;

DROP INDEX IF EXISTS public.ix_skill__title;
DROP INDEX IF EXISTS public.ix_skill__important;

DROP INDEX IF EXISTS public.ix_vacancy__created;
DROP INDEX IF EXISTS public.ix_vacancy__ldepartment;
DROP INDEX IF EXISTS public.ix_vacancy__area;
DROP INDEX IF EXISTS public.ix_vacancy__status;
DROP INDEX IF EXISTS public.ix_vacancy__title;
DROP INDEX IF EXISTS public.ix_vacancy__templated_id;

DROP TABLE IF EXISTS public.vacancy;

DROP TYPE IF EXISTS public.STATUS;

DROP SCHEMA IF EXISTS vacancy;