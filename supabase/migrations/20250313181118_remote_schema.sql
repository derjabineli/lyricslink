revoke delete on table "public"."songs" from "anon";

revoke insert on table "public"."songs" from "anon";

revoke references on table "public"."songs" from "anon";

revoke select on table "public"."songs" from "anon";

revoke trigger on table "public"."songs" from "anon";

revoke truncate on table "public"."songs" from "anon";

revoke update on table "public"."songs" from "anon";

revoke delete on table "public"."songs" from "authenticated";

revoke insert on table "public"."songs" from "authenticated";

revoke references on table "public"."songs" from "authenticated";

revoke select on table "public"."songs" from "authenticated";

revoke trigger on table "public"."songs" from "authenticated";

revoke truncate on table "public"."songs" from "authenticated";

revoke update on table "public"."songs" from "authenticated";

revoke delete on table "public"."songs" from "service_role";

revoke insert on table "public"."songs" from "service_role";

revoke references on table "public"."songs" from "service_role";

revoke select on table "public"."songs" from "service_role";

revoke trigger on table "public"."songs" from "service_role";

revoke truncate on table "public"."songs" from "service_role";

revoke update on table "public"."songs" from "service_role";

alter table "public"."songs" drop constraint "songs_pkey";

drop index if exists "public"."songs_pkey";

drop table "public"."songs";


