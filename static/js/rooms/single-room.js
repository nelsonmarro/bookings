const myButton = document.getElementById("error-btn"); // Or use querySelector, etc.
if (myButton) {
  myButton.click();
}

const form = document.getElementById("form-availability");
let formData = new FormData(form);
form.addEventListener("submit", function (event) {
  event.preventDefault();
  fetch("/reservation-json", {
    method: "POST",
    body: formData,
  })
    .then((response) => response.json())
    .then((data) => {
      console.log(data);
    });
});
