document.addEventListener("alpine:init", () => {
  Alpine.data("registerComponent", () => ({
    firstName: "",
    lastName: "",
    email: "",
    password: "",
    errorMessage: "",

    async register() {
      validPassword = validatePassword(this.password)
      if (!validPassword.valid) {
        this.errorMessage = validPassword.message
        return
      }

      try {
        let response = await fetch("/api/register", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            firstName: this.firstName,
            lastName: this.lastName,
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
          window.location.href = data.redirect
        }
      } catch (error) {
        console.error(error)
        this.errorMessage = "An error occurred. Please try again."
      }
    },
  }))
})

function validatePassword(password) {
  const regex = new RegExp("^(?=.*[a-z])(?=.*[A-Z]).{8,}$")

  return {
    valid: regex.test(password),
    message:
      "Password must be at least 8 characters long, include at least one uppercase letter, and at least one lowercase letter",
  }
}
