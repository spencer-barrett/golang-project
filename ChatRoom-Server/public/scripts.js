const chatLog = document.getElementById('chat-log');
const messageForm = document.getElementById('message-form');
const messageInput = document.getElementById('message-input');

// Ask for the user's name on page load
const author = prompt("Please enter your name:");

// Establish a WebSocket connection
const socket = new WebSocket("ws://" + window.location.host + "/ws");

// Event listener for incoming messages
socket.addEventListener("message", function(event) {
    const message = JSON.parse(event.data);
    const p = document.createElement('p');
    p.innerHTML = `<span class="author">${message.author}:</span> ${message.content}`;
    chatLog.appendChild(p);
    chatLog.scrollTop = chatLog.scrollHeight; // Auto-scroll to the bottom
});

// Handle form submission
messageForm.addEventListener("submit", function(event) {
    event.preventDefault();
    const messageContent = messageInput.value;
    if (messageContent && author) {
        const messageToSend = {
            author: author,
            content: messageContent
        };
        socket.send(JSON.stringify(messageToSend));
        messageInput.value = "";
    }
});