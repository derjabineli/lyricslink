function showArrangements(element) {
  arrangements = element.querySelector(".arrangements")
  if (arrangements.classList.contains("hidden")) {
    arrangements.classList.remove("hidden")
    element.classList.remove("hover:bg-gray-100")
  } else {
    arrangements.classList.add("hidden")
    element.classList.add("hover:bg-gray-100")
  }
}

function showAddSongModal() {
  songModal = document.getElementById("newSongModal")
  if (songModal.classList.contains("hidden")) {
    songModal.classList.remove("hidden")
  } else {
    songModal.classList.add("hidden")
  }
}

document.addEventListener("DOMContentLoaded", function () {
  document.querySelectorAll(".arrangements").forEach((arrangement) => {
    arrangement.addEventListener("click", function (event) {
      event.stopPropagation() // Prevents the event from bubbling up
    })
  })
})
