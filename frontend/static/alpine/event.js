document.addEventListener("alpine:init", () => {
  Alpine.store("event", {
    songModalOpen: false,
  })
  Alpine.data("eventComponent", () => ({
    id: null,
    name: "",
    date: "",

    async updateEvent() {
      try {
        let response = await fetch("/api/events", {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            id: this.id,
            date: this.date,
          }),
        })

        console.log(response)
      } catch (error) {
        console.log(error)
      }
    },
    async changeArrangement(arrangement_id) {},
  }))
})
