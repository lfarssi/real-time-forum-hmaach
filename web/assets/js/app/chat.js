import { getConvertation } from './api.js';
import { formatTime, showErrorPage, showNotification, debounce } from './utils.js';
import { sendMessage } from './websocket.js';

let currentPage = 1;
let isLoadingMessages = false;
let hasMoreMessages = true;
export let chatID = null;

export const showDirectMessages = async (id) => {
    currentPage = 1;
    isLoadingMessages = false;
    hasMoreMessages = true;

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

            <div class="messages-container">
                <div class="loading-indicator">
                    <i class="fa-solid fa-spinner fa-spin"></i> Loading more messages...
                </div>
            </div>


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
        setupMessageScroll();
    }
};

const setupMessageForm = () => {
    const messageForm = document.getElementById('message-form');
    const messagesContainer = document.querySelector('.messages-container');

    messageForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const input = messageForm.querySelector('input');
        const message = input.value.trim();

        if (message && chatID) {
            // Send message via WebSocket
            sendMessage(chatID, message);

            currentPage = 1;
            isLoadingMessages = false;
            hasMoreMessages = true;
            input.value = '';

            const loadingIndicator = messagesContainer.querySelector('.loading-indicator');
            messagesContainer.innerHTML = '';
            messagesContainer.appendChild(loadingIndicator);
            await loadConversation();
        } else {
            if (!chatID) showNotification('error', 'something wrong happend, please try later')
            if (!message) showNotification('error', 'please write a message first')
        }
    });
};

export const loadConversation = async () => {
    try {
        const token = localStorage.getItem('token');
        const response = await getConvertation(chatID, currentPage, token);
        if (response.status !== 200) throw response

        if (response.sender) {
            // Update recipient info
            document.getElementById('recipient-name').textContent = response.sender.nickname;
            document.getElementById('recipient-avatar').src = `https://ui-avatars.com/api/?name=${response.sender.nickname}`;
            renderMessages(response.messages);
        }
    } catch (error) {
        showErrorPage(error.status, error.message);
    }
};

const renderMessages = (messages) => {
    const messagesContainer = document.querySelector('.messages-container');
    const loadingIndicator = messagesContainer.querySelector('.loading-indicator');
    const userId = JSON.parse(localStorage.getItem('user')).id;

    if (!messages || messages.length === 0) {
        loadingIndicator.style.display = 'none';
        hasMoreMessages = false;
        return;
    }

    // Save scroll position
    const isAtBottom = messagesContainer.scrollHeight - messagesContainer.scrollTop === messagesContainer.clientHeight;

    messages.forEach(message => {
        const messageDiv = document.createElement('div');
        messageDiv.className = `message ${message.sender_id === userId ? 'sent' : 'received'}`;
        messageDiv.innerHTML = `
            <div class="message-content">
                <pre>${message.content}</pre>
                <span class="timestamp">${formatTime(message.sent_at)}</span>
            </div>
        `;

        // Insert after loading indicator
        loadingIndicator.after(messageDiv);
    });

    if (messages.length < 10) {
        loadingIndicator.style.display = 'none';
        hasMoreMessages = false;
    } else {
        loadingIndicator.style.display = 'block';
    }

    // Maintain scroll position
    if (currentPage === 1 || isAtBottom) {
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }
};

export const appendMessage = async () => {
    currentPage = 1;
    isLoadingMessages = false;
    hasMoreMessages = true;


    // Clear messages container except loading indicator
    const messagesContainer = document.querySelector('.messages-container');
    const loadingIndicator = messagesContainer.querySelector('.loading-indicator');
    messagesContainer.innerHTML = '';
    messagesContainer.appendChild(loadingIndicator);

    // Reload conversation
    await loadConversation();
};

const setupMessageScroll = () => {
    const messagesContainer = document.querySelector('.messages-container');

    const scrollFunc = debounce(async () => {
        if (isLoadingMessages || !hasMoreMessages) return;

        // Check if scrolled near top (100px threshold)
        if (messagesContainer.scrollTop <= 100) {
            isLoadingMessages = true;

            try {
                currentPage++;
                const token = localStorage.getItem('token');
                const response = await getConvertation(chatID, currentPage, token);

                if (response.status !== 200) throw response;

                // Save scroll position before adding new messages
                const oldScrollHeight = messagesContainer.scrollHeight;

                renderMessages(response.messages);

                // Adjust scroll position to avoid jumping
                if (hasMoreMessages) {
                    const newScrollHeight = messagesContainer.scrollHeight;
                    messagesContainer.scrollTop = newScrollHeight - oldScrollHeight;
                }
            } catch (error) {
                showErrorPage(error.status, error.message);
            } finally {
                isLoadingMessages = false;
            }
        }
    }, 800);

    messagesContainer.removeEventListener('scroll', scrollFunc);
    messagesContainer.addEventListener('scroll', scrollFunc);
};
