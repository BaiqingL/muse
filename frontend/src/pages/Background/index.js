console.log('This is the best page.');
console.log('Put the background scripts here.');

chrome.runtime.onMessage.addListener(function (request, sender, sendResponse) {
  fetch('http://localhost:8080/api/iterate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      html: request.html,
      prompt: request.prompt,
    }),
  });
});

console.log('listener added');
