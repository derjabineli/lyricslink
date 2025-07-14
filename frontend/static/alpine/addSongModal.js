document.addEventListener("alpine:init", () => {
  console.log("MOUNTED");
  Alpine.data("newSongModalComponent", () => ({
    query: "",
    songs: {},
    arrangements: [],

    async querySongs() {
      if (this.query.length === 0) {
        this.songs = [];
        return;
      }

      try {
        let response = await fetch("/api/songs", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            query: this.query,
          }),
        });

        let data = await response.json();
        if (!response.ok) {
          this.errorMessage =
            data.error || "Something went wrong. Please try again";
        } else {
          this.songs = data;
        }
      } catch (error) {
        console.error(error);
        this.errorMessage = "An error occurred. Please try again.";
      }
    },

    async getArrangements(key) {
      try {
        let response = await fetch(`/api/songs/${key}`, {
          method: "GET",
        });
        let data = await response.json();
        if (response.ok) {
          songSeachDiv = document.getElementById("songSearch");
          arrangementFormDiv = document.getElementById("arrangementForm");
          songSeachDiv.classList.add("hidden");
          arrangementFormDiv.classList.remove("hidden");
          this.arrangements = data;
        }
      } catch (error) {}
    },

    async addArrangementToEvent(arrangement_id, song_id) {
      pathname = window.location.pathname;
      event_id = pathname.split("/events/").join("");
      try {
        let response = await fetch(`/api/events_arrangements`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            event_id: event_id,
            arrangement_id: arrangement_id,
            song_id: song_id,
          }),
        });

        if (response.ok) {
          window.location.reload();
        }
      } catch (error) {}
    },
  }));
});

function checkIfSongInEvent(songs, key) {
  return songs.some((song) => key in song);
}
