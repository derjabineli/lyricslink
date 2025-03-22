document.addEventListener("alpine:init", () => {
  Alpine.data("userSettings", () => ({
    errorMessage: "",
    async syncPCSongs() {
      try {
        let response = await fetch("/api/")
      } catch (error) {
        console.log
      }
    },
    async openLoginPopup() {
      const loginWindow = window.open(
        "https://api.planningcenteronline.com/oauth/authorize?client_id=a880698546e79438bc6ebd4e0df5f4cf94a0b434b8cf7378d526984230183bf9&redirect_uri=http://localhost:3005/pc/callback&response_type=code&scope=services people",
        "oauth_popup",
        "width=600,height=400"
      )

      if (!loginWindow) {
        console.error("Popup blocked! Allow popups for this site.")
        return
      }

      const interval = setInterval(() => {
        if (loginWindow.closed) {
          clearInterval(interval)
          console.log("Login popup closed, check authentication status.")
          // Make a request to check session or token
        }
      }, 1000)
    },
  }))
})
