import { initMobileMenu } from "./mobileMenu.js";
import { setupDropdown } from "./dropdown.js";
import { setupGlobalEventListenerForDropdowns } from "./utils.js";

// Initialize the mobile menu functionality
initMobileMenu();

// Initialize User Profile Dropdown
const userProfileDropdownControl = setupDropdown(
  "user-menu-button",
  "user-menu-dropdown",
);

// Initialize Rooms Dropdown
const roomsDropdownControl = setupDropdown(
  "rooms-menu-button",
  "rooms-menu-dropdown",
);

setupGlobalEventListenerForDropdowns(
  userProfileDropdownControl,
  roomsDropdownControl,
);
