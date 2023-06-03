var highlightedElement = null; // Keep track of the currently highlighted element

function highlightElement(event) {
  var x = event.clientX; // Get the X coordinate of the cursor
  var y = event.clientY; // Get the Y coordinate of the cursor

  var element = document.elementFromPoint(x, y); // Get the element at the specified coordinates

  if (highlightedElement !== null) {
    highlightedElement.classList.remove('highlight'); // Remove the 'highlight' class from the previously highlighted element
  }

  element.classList.add('highlight'); // Add the 'highlight' class to the current element
  highlightedElement = element; // Update the highlighted element
}

async function replaceWithLoading() {
  if (highlightedElement !== null) {
    // Set position: relative to the highlighted element
    highlightedElement.style.position = 'relative';

    var loadingOverlay = document.createElement('div');
    loadingOverlay.className = 'loading-overlay';
    loadingOverlay.innerHTML = '<div class="loading">Loading...</div>';

    highlightedElement.appendChild(loadingOverlay); // Add the loading overlay to the highlighted element
  }
}

// Attach the event listener to track cursor movement
document.addEventListener('mousemove', highlightElement);
document.addEventListener('click', replaceWithLoading);
