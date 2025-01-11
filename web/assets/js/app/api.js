// Public Routes
export const registerUser = async (userData) => {
    const response = await fetch(`api/register`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(userData),
    });
    return response.json();
};

export const loginUser = async (credentials) => {
    const response = await fetch(`api/login`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(credentials),
    });
    return response.json();
};

// Authenticated Routes
export const getPosts = async (page = 1, token) => {
    const response = await fetch(`api/posts?page=${page}`, {
        headers: {
            'Authorization': `Bearer ${token}`,
        }
    })
    return response.json();
};

export const getComments = async (postId, page = 1, token) => {
    const response = await fetch(`api/posts/${postId}/comments?page=${page}`, {
        headers: {
            'Authorization': `Bearer ${token}`,
        }
    });
    return response.json();
};

export const createPost = async (postData, token) => {
    const response = await fetch(`api/posts/create`, {
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
    const response = await fetch(`api/comments/create`, {
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
    const response = await fetch(`api/logout`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`,
        },
    });
    return response.json();
};
