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
    return response.json();
};

export const getUsers = async (token) => {
    const response = await fetch(`api/users`, {
        headers: {
            'Authorization': `Bearer ${token}`,
        }
    })
    return response.json();
};

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

export const reactToPost = async (reactionData, token) => {
    const response = await fetch(`api/posts/react`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(reactionData),
    });
    return response.json();
};

export const getConvertation = async (senderId, token) => {
    const response = await fetch(`api/conversation/${senderId}`, {
        headers: {
            'Authorization': `Bearer ${token}`,
        }
    });
    return response.json();
};