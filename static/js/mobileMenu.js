export function initMobileMenu() {
  // Select the button that controls the mobile menu.
  // Note: The original selector was 'button[aria-controls="mobile-menu"]'.
  // For this demo, I'm using a more direct ID for simplicity, assuming it's unique.
  const mobileMenuButton = document.querySelector(
    'button[aria-controls="mobile-menu"]',
  );

  const mobileMenu = document.getElementById("mobile-menu");

  if (mobileMenuButton && mobileMenu) {
    // Find the open and close icons within the button.
    const mobileMenuOpenIcon = mobileMenuButton.querySelector("svg.block");
    const mobileMenuCloseIcon = mobileMenuButton.querySelector("svg.hidden");

    if (mobileMenuOpenIcon && mobileMenuCloseIcon) {
      mobileMenuButton.addEventListener("click", () => {
        // Check the current expanded state.
        const isExpanded =
          mobileMenuButton.getAttribute("aria-expanded") === "true" || false;

        // Toggle the ARIA attribute.
        mobileMenuButton.setAttribute("aria-expanded", !isExpanded);

        // Toggle visibility of the open and close icons.
        mobileMenuOpenIcon.classList.toggle("hidden");
        mobileMenuOpenIcon.classList.toggle("block");
        mobileMenuCloseIcon.classList.toggle("hidden");
        mobileMenuCloseIcon.classList.toggle("block");

        // Toggle visibility of the mobile menu itself.
        mobileMenu.classList.toggle("hidden");
      });
    } else {
      console.warn(
        "Mobile menu icons not found. Ensure 'svg.block' and 'svg.hidden' exist within the button.",
      );
    }
  } else {
    console.warn(
      "Mobile menu button or menu element not found. Check IDs: 'mobile-menu-toggle-button', 'mobile-menu'.",
    );
  }
}
