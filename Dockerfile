# Use a lightweight debian os
# as the base image
FROM debian:stable-slim

# COPY source destination
COPY lyriclink /bin/lyriclink
COPY frontend /frontend
COPY internal /internal
COPY sql /sql
COPY supabase /supabase

CMD ["/bin/lyriclink"]