document.addEventListener("alpine:init", () => {
  console.log("MOUNTED")
  Alpine.data("loginComponent", () => ({
    email: "",
    password: "",
    errorMessage: "",

    async login() {
      try {
        let response = await fetch("/api/login", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            email: this.email,
            password: this.password,
          }),
        })

        let data = await response.json()
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
