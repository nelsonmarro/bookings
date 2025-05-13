export function setupGlobalEventListenerForDropdowns(
  userProfileDropdownControl,
  roomsDropdownControl,
) {
  // --- Global Event Listeners to close dropdowns ---

  // Close dropdowns if a click happens outside of them
  document.addEventListener("click", (event) => {
    // Close User Profile Dropdown if open and click is outside
    if (
      userProfileDropdownControl &&
      userProfileDropdownControl.button.getAttribute("aria-expanded") === "true"
    ) {
      if (
        !userProfileDropdownControl.button.contains(event.target) &&
        !userProfileDropdownControl.dropdown.contains(event.target)
      ) {
        userProfileDropdownControl.hideFunc();
      }
    }

    // Close Rooms Dropdown if open and click is outside
    if (
      roomsDropdownControl &&
      roomsDropdownControl.button.getAttribute("aria-expanded") === "true"
    ) {
      if (
        !roomsDropdownControl.button.contains(event.target) &&
        !roomsDropdownControl.dropdown.contains(event.target)
      ) {
        roomsDropdownControl.hideFunc();
      }
    }
    // Add similar logic for other dropdowns if you have more
  });

  // Close dropdowns if the Escape key is pressed
  document.addEventListener("keydown", (event) => {
    if (event.key === "Escape") {
      // Close User Profile Dropdown
      if (
        userProfileDropdownControl &&
        userProfileDropdownControl.button.getAttribute("aria-expanded") ===
          "true"
      ) {
        userProfileDropdownControl.hideFunc(true); // Hide immediately
        userProfileDropdownControl.button.focus(); // Return focus to the button
      }

      // Close Rooms Dropdown
      if (
        roomsDropdownControl &&
        roomsDropdownControl.button.getAttribute("aria-expanded") === "true"
      ) {
        roomsDropdownControl.hideFunc(true); // Hide immediately
        roomsDropdownControl.button.focus(); // Return focus to the button
      }
      // Add similar logic for other dropdowns
    }
  });
}
