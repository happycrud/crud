CREATE TABLE "public"."user" (
    "id" serial NOT NULL PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "age" int4 NOT NULL,
    "ctime" timestamp(6) NOT NULL DEFAULT now(),
    "mtime" timestamp(6) NOT NULL DEFAULT now()
);