var isPickingEnabled = false;

var highlightedElement = null; // Keep track of the currently highlighted element
var selectedElement = null; // Keep track of the currently selected element

var popup = document.createElement('div');
popup.className = 'muse-popup';
popup.innerHTML = `
<p>What should we change?</p>
<textarea placeholder="Move this button to the top left and make it blue"></textarea>
<button>Submit</button>
`;

var toggle = document.createElement('button');
toggle.className = 'muse-toggle';
toggle.innerHTML = `
<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-brush-fill" viewBox="0 0 16 16">
  <path d="M15.825.12a.5.5 0 0 1 .132.584c-1.53 3.43-4.743 8.17-7.095 10.64a6.067 6.067 0 0 1-2.373 1.534c-.018.227-.06.538-.16.868-.201.659-.667 1.479-1.708 1.74a8.118 8.118 0 0 1-3.078.132 3.659 3.659 0 0 1-.562-.135 1.382 1.382 0 0 1-.466-.247.714.714 0 0 1-.204-.288.622.622 0 0 1 .004-.443c.095-.245.316-.38.461-.452.394-.197.625-.453.867-.826.095-.144.184-.297.287-.472l.117-.198c.151-.255.326-.54.546-.848.528-.739 1.201-.925 1.746-.896.126.007.243.025.348.048.062-.172.142-.38.238-.608.261-.619.658-1.419 1.187-2.069 2.176-2.67 6.18-6.206 9.117-8.104a.5.5 0 0 1 .596.04z"/>
</svg>`;

toggle.addEventListener('click', function () {
  isPickingEnabled = !isPickingEnabled;

  if (isPickingEnabled) {
    enablePicker();
  } else {
    disablePicker();
  }

  console.log('isPickingEnabled: ' + isPickingEnabled);
});

document.body.appendChild(toggle); // Add the toggle to the DOM
document.body.appendChild(popup); // Add the popup to the DOM

function disablePicker() {
  document.removeEventListener('mousemove', highlightElement);
  document.removeEventListener('click', replaceWithLoading);
  popup.style.display = 'none';

  if (highlightedElement !== null) {
    highlightedElement.classList.remove('muse-highlight');
  }

  if (selectedElement !== null) {
    selectedElement.classList.remove('muse-selected');
  }
}

function enablePicker() {
  document.addEventListener('mousemove', highlightElement);
  document.addEventListener('click', replaceWithLoading);
}

var textArea = popup.querySelector('textarea');
var button = popup.querySelector('button');

function isMuseElement(element) {
  return (
    element === popup ||
    popup.contains(element) ||
    element === toggle ||
    toggle.contains(element)
  );
}

button.addEventListener('click', function () {
  // remove any muse- classes from element
  var elementClasses = selectedElement.classList;

  for (var i = elementClasses.length - 1; i >= 0; i--) {
    var className = elementClasses[i];
    if (className.startsWith('muse-')) {
      elementClasses.remove(className);
    }
  }

  const html = selectedElement.outerHTML;

  chrome.runtime.sendMessage({
    html: html,
    prompt: textArea.value,
  });
});

function highlightElement(event) {
  var x = event.clientX; // Get the X coordinate of the cursor
  var y = event.clientY; // Get the Y coordinate of the cursor

  var element = document.elementFromPoint(x, y); // Get the element at the specified coordinates

  if (element === null || isMuseElement(element)) {
    // If the cursor is over the popup, don't highlight it
    return;
  }

  if (highlightedElement !== null) {
    highlightedElement.classList.remove('muse-highlight'); // Remove the 'highlight' class from the previously highlighted element
  }

  element.classList.add('muse-highlight'); // Add the 'highlight' class to the current element
  highlightedElement = element; // Update the highlighted element
}

async function replaceWithLoading(event) {
  // If clicked on the popup, don't do anything
  const element = event.target;
  if (element === null || isMuseElement(element)) {
    return;
  }

  if (highlightedElement !== null) {
    // Move the popup to the cursor position
    popup.style.left = event.pageX + 'px';
    popup.style.top = event.pageY + 'px';

    popup.style.display = 'block'; // Show the popup

    if (selectedElement !== null) {
      selectedElement.classList.remove('muse-selected'); // Remove the 'selected' class from the previously selected element
    }

    highlightedElement.classList.add('muse-selected'); // Add the 'selected' class to the current element
    selectedElement = highlightedElement; // Update the selected element
  }
}

disablePicker();
