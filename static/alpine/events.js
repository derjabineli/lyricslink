document.addEventListener("alpine:init", () => {
  console.log("MOUNTED")
  Alpine.data("modalComponent", () => ({
    name: "",
    date: "",

    async addEvent() {
      try {
        let response = await fetch("/api/events", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            name: this.name,
            date: this.date,
          }),
        })

        let data = await response.json()
        console.log(data)
        if (!response.ok) {
          this.errorMessage =
            data.error || "Something went wrong. Please try again"
        } else {
          this.errorMessage = ""
          window.location.href = "/dashboard"
        }
      } catch (error) {
        console.error(error)
        this.errorMessage = "An error occurred. Please try again."
      }
    },
  }))
})
