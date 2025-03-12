const modal = document.getElementById("eventModal")
const openModalBtn = document.getElementById("openModalBtn")
const closeModalBtn = document.getElementById("closeModalBtn")
const entryForm = document.getElementById("entryForm")

// Show modal with blur effect
openModalBtn.addEventListener("click", () => {
  modal.classList.remove("hidden")
})

// Hide modal
closeModalBtn.addEventListener("click", () => {
  modal.classList.add("hidden")
})

// Handle form submission
entryForm.addEventListener("submit", (e) => {
  e.preventDefault()
  alert("Form submitted!") // Replace with actual logic
  modal.classList.add("hidden")
})
