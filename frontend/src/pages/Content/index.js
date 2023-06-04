var highlightedElement = null; // Keep track of the currently highlighted element
var selectedElement = null; // Keep track of the currently selected element

var popup = document.createElement('div');
popup.className = 'muse-popup';
popup.innerHTML = `
<p>What should we change?</p>
<textarea placeholder="Move this button to the top left and make it blue"></textarea>
<button>Submit</button>
`;

var textArea = popup.querySelector('textarea');
var button = popup.querySelector('button');

button.addEventListener('click', function () {
  console.log(textArea.value);

  // get some defining characteristics of the element, like its tag name, class names, and text content
  var tagName = selectedElement.tagName;
  var classNames = Array.from(selectedElement.classList).filter(
    (className) => !className.startsWith('muse-')
  );
  var textContent = selectedElement.textContent;

  console.log(`tagName: ${tagName}`);
  console.log(`classNames: ${classNames} ${classNames.length}`);
  console.log(`textContent: ${textContent}`);
});

document.body.appendChild(popup); // Add the popup to the DOM

function highlightElement(event) {
  var x = event.clientX; // Get the X coordinate of the cursor
  var y = event.clientY; // Get the Y coordinate of the cursor

  var element = document.elementFromPoint(x, y); // Get the element at the specified coordinates

  if (element === popup || popup.contains(element) || element === null) {
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
  if (event.target === popup || popup.contains(event.target)) {
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

async function downloadZip() {
  const response = await fetch('/export', { // Change '/export' to your Go server's route
    method: 'GET'
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  const blob = await response.blob();

  const url = window.URL.createObjectURL(blob);

  const link = document.createElement('a');
  link.href = url;
  link.download = 'export.zip';

  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
}

// Attach the event listener to track cursor movement
document.addEventListener('mousemove', highlightElement);
document.addEventListener('click', replaceWithLoading);

// Attach the download function to the 'Export' button
document.querySelector('#exportButton').addEventListener('click', downloadZip);