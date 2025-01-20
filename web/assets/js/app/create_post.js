// create_post.js
import { createPost } from './api.js';
import { showFeed } from './feed.js';
import { showErrorPage, getFormData, showNotification } from './utils.js';

export const showCreatePost = () => {
    const mainContainer = document.querySelector('main')
    mainContainer.innerHTML = /*html*/`
    <form id="new-post-form" class="create-post-form">
        <h2>Create New Post</h2>
        <div class="form-section">
            <label for="post-title">Title</label>
            <input type="text" id="post-title" name="title" placeholder="Enter your post title..." required>
        </div>

        <div class="form-section">
            <label>Categories</label>
            <div class="categories-container">
                <div class="category" data-id="1" data-selected="false">
                    <i class="fa-solid fa-microchip"></i>
                    <span>Technology</span>
                </div>
                <div class="category" data-id="2" data-selected="false">
                    <i class="fa-solid fa-futbol"></i>
                    <span>Sport</span>
                </div>
                <div class="category" data-id="3" data-selected="false">
                    <i class="fa-solid fa-chart-line"></i>
                    <span>Business</span>
                </div>
                <div class="category" data-id="4" data-selected="false">
                    <i class="fa-solid fa-heart-pulse"></i>
                    <span>Health</span>
                </div>
                <div class="category" data-id="5" data-selected="false">
                    <i class="fa-solid fa-newspaper"></i>
                    <span>News</span>
                </div>
            </div>
        </div>

        <div class="form-section">
            <label for="post-content">Content</label>
            <textarea id="post-content" name="content" placeholder="Write your post content here..." required></textarea>
        </div>

        <div class="form-actions">
            <button type="submit" class="submit-btn">Create Post</button>
        </div>
    </form>
    `

    // Initialize components
    setupCategorySelection();
    setupFormSubmission();
};

const setupCategorySelection = () => {
    const categories = document.querySelectorAll('.category');
    categories.forEach(category => {
        category.addEventListener('click', () => {
            const isSelected = category.getAttribute('data-selected') === 'true';
            category.setAttribute('data-selected', !isSelected);
        });
    });
};

const setupFormSubmission = () => {
    const form = document.getElementById('new-post-form');
    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        // Get selected categories
        const selectedCategories = Array.from(document.querySelectorAll('.category[data-selected="true"]'))
            .map(cat => parseInt(cat.getAttribute('data-id')));
        
        if (selectedCategories.length === 0) {
            showNotification('error', 'Please select at least one category');
            return;
        }

        const formData = getFormData(form);
        const postData = { ...formData, categories: selectedCategories };

        try {
            const token = localStorage.getItem('token');
            const response = await createPost(postData, token);

            if (response.status !== 200) throw response
            showFeed();
        } catch (error) {
            showErrorPage(error.status, error.message)
        }
    });
};