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
    async changeArrangement(
      event_arrangement_id,
      new_arrangement_id,
      event_id
    ) {
      try {
        let response = await fetch(
          `/api/events_arrangements/${event_arrangement_id}`,
          {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
              arrangement_id: new_arrangement_id,
              event_id: event_id,
            }),
          }
        )

        console.log(response)
        if (response.ok) {
          window.location.reload()
        }
      } catch (error) {
        console.log(error)
      }
    },
    async deleteEvent(event_id) {
      try {
        let response = await fetch(`/api/events/${event_id}`, {
          method: "DELETE",
        })

        if (response.status == 204) {
          window.location.href = "/dashboard"
        }
      } catch (error) {
        console.log(error)
      }
    },
    async deleteEventArrangement(event_arrangement_id) {
      try {
        let response = await fetch(
          `/api/events_arrangements/${event_arrangement_id}`,
          {
            method: "DELETE",
          }
        )

        if (response.status == 204) {
          window.location.reload()
        }
      } catch (error) {
        console.log(error)
      }
    },
  }))
})
