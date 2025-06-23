addEventListener("DOMContentLoaded", function () {
  const myButton = document.getElementById("error-btn");
  if (myButton) {
    myButton.click();
  }

  const form = document.getElementById("form-availability");

  form.addEventListener("submit", function (event) {
    event.preventDefault();
    let formData = new FormData(form);
    formData.append("room_id", "1"); // Make sure this room_id is correct or dynamic if needed

    fetch("/reservation-json", {
      method: "POST",
      body: formData,
    })
      .then(async (response) => {
        if (!response.ok) {
          const errData = await response.json();
          return await Promise.reject(errData);
        }
        return response.json();
      })
      .then((data) => {
        if (data.ok) {
          Swal.fire({
            icon: "success",
            title: "Room Available!",
            html: `
              <div class="py-3">
                 <a href="/book-room?id=${data.room_id}&s=${data.start_date}&e=${data.end_date}"
                 class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Book Now!</a>
              </div>
                  `,
            showConfirmButton: false,
            showCancelButton: true,
          });
        } else {
          Swal.fire({
            icon: "error",
            title: "Not Available",
            text:
              data.message ||
              "Sorry, this room is not available for your dates.", // Use backend message if available
            confirmButtonText: "Okay",
          });
        }
      })
      .catch((error) => {
        console.error("Fetch error:", error);
        let errorMessage = "Could not check availability. Please try again.";
        if (error && error.message) {
          // If error is an object with a message property
          errorMessage = error.message;
        } else if (typeof error === "string") {
          // If error is a simple string
          errorMessage = error;
        }

        Swal.fire({
          icon: "error",
          title: "Oops...",
          text: errorMessage,
          confirmButtonText: "Close",
        });
      });
  });
});
