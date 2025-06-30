async function logout() {
  let response = await fetch("/api/logout");
  if (response.status == 204) {
    window.location.href = "/login";
  }
}
