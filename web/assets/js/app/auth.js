import {
    registerUser,
    loginUser,
    logoutUser
} from './api.js'
import { writeMessage } from './utils.js';
import { showFeed } from './feed.js';


export const showAuth = () => {
    const authContainer = document.getElementById("auth-container");
    const feed = document.getElementById("feed");
    authContainer.style.display = "flex";
    feed.style.display = "none";
};

export const handleRegistration = async () => {
    const registerForm = document.getElementById("register-submit");
    if (registerForm) {
        registerForm.addEventListener("click", async (e) => {
            const user = {
                "first_name": document.getElementById("register-first-name")?.value,
                "last_name": document.getElementById("register-last-name")?.value,
                "email": document.getElementById("register-identifier").value,
                "nickname": document.getElementById("register-nickname")?.value,
                "gender": document.getElementById("register-gender")?.value,
                "age": parseInt(document.getElementById("register-age")?.value, 10),
                "password": document.getElementById("register-password").value,
                "password_confirmation": document.getElementById("register-password").value,
            };

            try {
                const response = await registerUser(user);
                if (response.ok) {
                    writeMessage("register-error", "Registration successful")
                } else {
                    throw response.message;
                }
            } catch (error) {
                writeMessage("register-error", error)
            }
        });
    }
};


export const handleLogin = async () => {
    const loginForm = document.getElementById("login-submit");
    if (loginForm) {
        loginForm.addEventListener("click", async () => {
            const credentials = {
                "nickname": document.getElementById("login-identifier").value,
                "password": document.getElementById("login-password").value
            };
            try {
                const response = await loginUser(credentials);
                if (response.message === "success" && response.user && response.token) {
                    // Clear existing user  and token
                    localStorage.removeItem("user");
                    localStorage.removeItem("token");

                    // Save new user data and token
                    localStorage.setItem("user", JSON.stringify(response.user));
                    localStorage.setItem("token", response.token);

                    showFeed(response.user);
                } else {
                    throw response.message;
                }
            } catch (error) {
                writeMessage("login-error", error);
            }
        });
    }
};


export const handleLogout = () => {
    const logoutBtn = document.getElementById("logout-submit");
    if (logoutBtn) {
        logoutBtn.addEventListener("click", async () => {
            try {
                const token = localStorage.getItem("token");
                const response = await logoutUser(token);

                if (response.message === "success") {
                    localStorage.removeItem("user");
                    localStorage.removeItem("token");
                    showAuth();
                } else {
                    console.error(response.message);
                }
            } catch (error) {
                console.error(error);
            }
        });
    }
};