document.addEventListener("alpine:init", () => {
  console.log("MOUNTED")
  Alpine.data("newSongModalComponent", () => ({
    query: "",
    songs: [],

    async querySongs() {
      if (this.query.length === 0) {
        this.songs = [] // ✅ Clear results when input is empty
        return
      }

      try {
        let response = await fetch("/api/songs", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            query: this.query,
          }),
        })

        let data = await response.json()
        console.log(data)
        if (!response.ok) {
          this.errorMessage =
            data.error || "Something went wrong. Please try again"
        } else {
          this.songs = data
        }
      } catch (error) {
        console.error(error)
        this.errorMessage = "An error occurred. Please try again."
      }
    },
  }))
})
