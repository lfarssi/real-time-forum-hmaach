import { handleUnauthorized } from "./utils.js";

export const authUser = async (formData, path) => {
    const response = await fetch(path, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
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
    const data = await response.json();
    if (handleUnauthorized(data)) return null;
    return data;
};

export const getUsers = async (token) => {
    const response = await fetch(`api/users`, {
        headers: {
            'Authorization': `Bearer ${token}`,
        }
    })
    const data = await response.json();
    if (handleUnauthorized(data)) return null;
    return data;
};

export const getPosts = async (page = 1, token) => {
    const response = await fetch(`api/posts?page=${page}`, {
        headers: {
            'Authorization': `Bearer ${token}`,
        }
    })
    const data = await response.json();
    if (handleUnauthorized(data)) return null;
    return data;
};

export const getComments = async (postId, page = 1, token) => {
    const response = await fetch(`api/posts/${postId}/comments?page=${page}`, {
        headers: {
            'Authorization': `Bearer ${token}`,
        }
    });
    const data = await response.json();
    if (handleUnauthorized(data)) return null;
    return data;
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
    const data = await response.json();
    if (handleUnauthorized(data)) return null;
    return data;
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
    const data = await response.json();
    if (handleUnauthorized(data)) return null;
    return data;
};

export const reactToPost = async (reactionData, token) => {
    const response = await fetch(`api/posts/react`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(reactionData),
    });
    const data = await response.json();
    if (handleUnauthorized(data)) return null;
    return data;
};

export const getConvertation = async (senderId, page = 1, token) => {
    const response = await fetch(`api/conversation/${senderId}?page=${page}`, {
        headers: {
            'Authorization': `Bearer ${token}`,
        }
    });
    const data = await response.json();
    if (handleUnauthorized(data)) return null;
    return data;
};