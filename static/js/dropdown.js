/**
 * Sets up a dropdown menu with click-to-toggle functionality and animations.
 * @param {string} buttonId - The ID of the button that triggers the dropdown.
 * @param {string} dropdownId - The ID of the dropdown menu element.
 * @returns {object|null} An object containing references to the button, dropdown,
 * and a hide function, or null if elements are not found.
 */
export function setupDropdown(buttonId, dropdownId) {
  const button = document.getElementById(buttonId);
  const dropdown = document.getElementById(dropdownId);

  if (!button || !dropdown) {
    console.warn(
      `Dropdown elements not found for button ID: ${buttonId} or dropdown ID: ${dropdownId}`,
    );
    return null;
  }

  // Ensure base transition and transform classes are on the dropdown for animations
  // These classes are expected to be defined in your CSS (e.g., Tailwind CSS or custom)
  if (!dropdown.classList.contains("transition"))
    dropdown.classList.add("transition");
  if (!dropdown.classList.contains("transform"))
    dropdown.classList.add("transform");

  /**
   * Shows the dropdown menu with an animation.
   */
  const showDropdown = () => {
    button.setAttribute("aria-expanded", "true");
    dropdown.classList.remove("hidden"); // Make it visible first

    // Animation classes for showing (Tailwind-like)
    dropdown.classList.remove(
      "ease-in",
      "duration-75",
      "opacity-0",
      "scale-95",
    );
    dropdown.classList.add(
      "ease-out",
      "duration-100",
      "opacity-100",
      "scale-100",
    );

    // The requestAnimationFrame and class juggling is a common pattern to ensure
    // the "from" state of the animation is applied before the "to" state.
    // For simpler transitions, you might not need this level of detail.
    // However, to match the original logic's intent for Tailwind-like transitions:
    // Start from a hidden-like animated state
    dropdown.classList.add("opacity-0", "scale-95");
    requestAnimationFrame(() => {
      dropdown.classList.remove("opacity-0", "scale-95");
      // The 'opacity-100' and 'scale-100' should already be there from above,
      // but this ensures the transition plays correctly.
    });
  };

  /**
   * Hides the dropdown menu with an animation.
   * @param {boolean} [immediate=false] - If true, hides the dropdown immediately without animation.
   */
  const hideDropdown = (immediate = false) => {
    button.setAttribute("aria-expanded", "false");

    // Animation classes for hiding
    dropdown.classList.remove(
      "ease-out",
      "duration-100",
      "opacity-100",
      "scale-100",
    );
    dropdown.classList.add("ease-in", "duration-75", "opacity-0", "scale-95");

    const transitionDuration = immediate ? 0 : 75; // Match duration-75

    setTimeout(() => {
      // Only hide if it's still meant to be closed (e.g., user hasn't re-clicked)
      if (button.getAttribute("aria-expanded") === "false") {
        dropdown.classList.add("hidden");
      }
    }, transitionDuration);
  };

  // Event listener for the dropdown button
  button.addEventListener("click", (event) => {
    event.stopPropagation(); // Prevent click from bubbling up to document listener immediately
    const isExpanded = button.getAttribute("aria-expanded") === "true";
    if (isExpanded) {
      hideDropdown();
    } else {
      // Optional: Close other dropdowns before showing this one.
      // You would need a way to manage all active dropdowns if you implement this.
      // For now, this example allows multiple dropdowns to be open.
      // closeOtherDropdowns(dropdownId);
      showDropdown();
    }
  });

  // Return control object for external management (e.g., global close listeners)
  return { button, dropdown, hideFunc: hideDropdown, showFunc: showDropdown };
}
