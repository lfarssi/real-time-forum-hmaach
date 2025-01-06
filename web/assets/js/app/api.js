const BASE_URL = 'http://localhost:8080';

// Public Routes
export const getPosts = async (page = 1) => {
    const response = await fetch(`${BASE_URL}/posts?page=${page}`);
    return response.json();
};

export const getComments = async (postId, page = 1) => {
    const response = await fetch(`${BASE_URL}/posts/${postId}/comments?page=${page}`);
    return response.json();
};

export const registerUser = async (userData) => {
    const response = await fetch(`${BASE_URL}/register`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(userData),
    });
    return response.json();
};

export const loginUser = async (credentials) => {
    const response = await fetch(`${BASE_URL}/login`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(credentials),
    });
    return response.json();
};

// Authenticated Routes
export const createPost = async (postData, token) => {
    const response = await fetch(`${BASE_URL}/posts/create`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(postData),
    });
    return response.json();
};

export const createComment = async (commentData, token) => {
    const response = await fetch(`${BASE_URL}/comments/create`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(commentData),
    });
    return response.json();
};

export const logoutUser = async (token) => {
    const response = await fetch(`${BASE_URL}/logout`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`,
        },
    });
    return response.json();
};
