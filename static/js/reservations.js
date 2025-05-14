"use strict";

let forms = document.querySelectorAll("form");

Array.prototype.filter.call(forms, function (form) {
  form.addEventListener(
    "submit",
    function (event) {
      if (form.checkValidity() === false) {
        event.preventDefault();
        event.stopPropagation();
      }
    },
    false,
  );
});
