import { getConvertation } from './api.js';
import { formatTime, showErrorPage } from './utils.js';
import { sendMessage } from './websocket.js';

export let chatID = null;
export const showDirectMessages = async (id) => {
    const mainContainer = document.querySelector('main');
    mainContainer.innerHTML = `
        <div class="chat-main">
            <div class="chat-header">
                <div class="chat-user-info">
                    <img src="" alt="profile" id="recipient-avatar">
                    <div>
                        <span class="username" id="recipient-name"></span>
                    </div>
                </div>
            </div>

            <div class="messages-container"></div>

            <div class="message-input-container">
                <form id="message-form">
                    <input type="text" placeholder="Type your message..." required>
                    <button type="submit">
                        <i class="fa-solid fa-paper-plane"></i>
                    </button>
                </form>
            </div>
        </div>
    `;
    setupMessageForm();
    if (id) {
        chatID = id
        await loadConversation();
    }
};

const setupMessageForm = () => {
    const messageForm = document.getElementById('message-form');
    const messagesContainer = document.querySelector('.messages-container');

    messageForm.addEventListener('submit', (e) => {
        e.preventDefault();
        const input = messageForm.querySelector('input');
        const message = input.value.trim();

        if (message && chatID) {
            // Send message via WebSocket
            sendMessage(chatID, message);

            // Add message to UI
            appendMessage({content: message, sender_id: JSON.parse(localStorage.getItem('user')).id, sent_at: new Date().toISOString()});

            input.value = '';
            messagesContainer.scrollTop = messagesContainer.scrollHeight;
        }
    });
};

export const loadConversation = async () => {
    try {
        const token = localStorage.getItem('token');
        const response = await getConvertation(chatID, token);
        if (response.sender) {
            // Update recipient info
            document.getElementById('recipient-name').textContent = response.sender.nickname;
            document.getElementById('recipient-avatar').src = `https://ui-avatars.com/api/?name=${response.sender.nickname}`;
        }
        if (response.sender && response.messages) {
            // Clear and load messages
            const messagesContainer = document.querySelector('.messages-container');
            messagesContainer.innerHTML = '';
            
            response.messages.reverse().forEach(message => {
                appendMessage(message);
            });
            messagesContainer.scrollTop = messagesContainer.scrollHeight;
        }
    } catch (error) {
        showErrorPage(error);
    }
};

export const appendMessage = (message) => {
    const messagesContainer = document.querySelector('.messages-container');
    const userId = JSON.parse(localStorage.getItem('user')).id;
    
    const messageDiv = document.createElement('div');
    messageDiv.className = `message ${message.sender_id === userId ? 'sent' : 'received'}`;
    messageDiv.innerHTML = `
        <div class="message-content">
            <p>${message.content}</p>
            <span class="timestamp">${formatTime(message.sent_at)}</span>
        </div>
    `;
    
    messagesContainer.appendChild(messageDiv);
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
};