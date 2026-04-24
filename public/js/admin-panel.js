(function () {
  var sidebar = document.getElementById("adminSidebar");
  var overlay = document.getElementById("adminSidebarOverlay");
  var toggle = document.getElementById("adminSidebarToggle");
  var fullBtn = document.getElementById("adminFullscreenBtn");

  if (toggle && sidebar && overlay) {
    toggle.addEventListener("click", function () {
      sidebar.classList.toggle("show");
      overlay.classList.toggle("show");
    });
    overlay.addEventListener("click", function () {
      sidebar.classList.remove("show");
      overlay.classList.remove("show");
    });
  }

  if (fullBtn) {
    fullBtn.addEventListener("click", function () {
      if (!document.fullscreenElement) {
        document.documentElement.requestFullscreen().catch(function () {});
      } else {
        document.exitFullscreen().catch(function () {});
      }
    });
  }
})();
