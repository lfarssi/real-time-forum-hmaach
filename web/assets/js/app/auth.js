import { authUser, logoutUser } from './api.js'
import { showFeed } from './feed.js';
import { setupLayout } from './layout.js';
import { getFormData, showErrorPage } from './utils.js'
import { closeWebsocket } from './websocket.js';

export const showAuth = () => {
    document.body.innerHTML = ``;
    const formContainer = document.createElement('div');
    formContainer.id = 'form-container';
    formContainer.innerHTML = `
    <div class="container">
        <div class="register">
            <h2>Register</h2>
            <form id="register-form" action="/api/register" method="post">
                <div class="flex">
                    <input class="split" type="text" name="firstName" placeholder="First Name" pattern="[A-Za-z]+"
                        minlength="3" maxlength="20" required title="Please enter a valid first name (letters only)">
                    <input class="split" type="text" name="lastName" placeholder="Last Name" pattern="[A-Za-z]+"
                        minlength="3" maxlength="20" required title="Please enter a valid last name (letters only)">
                </div>
                <input type="text" name="nickname" placeholder="nickname" pattern="[a-z0-9]+" minlength="3"
                    maxlength="20" required title="nickname must be 3-20 characters">
                <div class="flex">
                    <input class="split" id="age-input" type="number" name="age" placeholder="Age" min="18"
                        max="120" required title="You must be at least 18 years old">
                    <select class="split" name="gender" required>
                        <option value="">Select Gender</option>
                        <option value="male">Male</option>
                        <option value="female">Female</option>
                    </select>
                </div>
                <input type="email" name="email" placeholder="Email" required>
                <input type="password" name="password" placeholder="Password" required>
                <p id="register-error" ></p>
                <input type="submit" value="Register">
            </form>
        </div>

        <div class="login">
            <h2>Login</h2>
            <form id="login-form" name="formLogin" action="/api/login" method="post">
                <input type="text" name="identifier" placeholder="Email or nickname" required>
                <input type="password" name="password" placeholder="Password" required>
                <p id="login-error" ></p>
                <input type="submit" value="Login">
            </form>
        </div>
    </div>
    `;
    document.body.appendChild(formContainer);
    setupFormToggle();
    FormSubmission();
};

const setupFormToggle = () => {
    const formContainer = document.querySelector("#form-container");
    const registerHeader = document.querySelector(".register h2");
    const loginHeader = document.querySelector(".login h2");

    // Toggle between login and register forms
    loginHeader.addEventListener("click", () => {
        formContainer.classList.add("active");
    });

    registerHeader.addEventListener("click", () => {
        formContainer.classList.remove("active");
    });
};

const FormSubmission = () => {
    const registerForm = document.querySelector('#register-form');
    const loginForm = document.querySelector('#login-form');
    const loginMsg = document.querySelector('#login-error');
    const registerMsg = document.querySelector('#register-error');

    [registerForm, loginForm].forEach(form => {
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            loginMsg.textContent = '';
            registerMsg.textContent = '';

            const formData = getFormData(form)
            try {
                const response = await authUser(formData, form.action);
                if (response.status === 200 && response.user && response.token) {
                    // Clear existing user and token
                    localStorage.removeItem("user");
                    localStorage.removeItem("token");
                    // Save new user data and token
                    localStorage.setItem("user", JSON.stringify(response.user));
                    localStorage.setItem("token", response.token);
                    setupLayout();
                    showFeed();
                } else {
                    throw response;
                }
            } catch (error) {
                if (form.id === 'register-form' && error.status === 400) {
                    registerMsg.textContent = error.message
                } else if (form.id === 'login-form') {
                    if (error.status === 500) {
                        showErrorPage(error)
                    } else {
                        loginMsg.textContent = error.message
                    }
                } else {
                    showErrorPage(error)
                }
            }
        });
    });
};

export const handleLogout = async () => {
    try {
        const token = localStorage.getItem('token');
        await logoutUser(token);
        localStorage.removeItem('user');
        localStorage.removeItem('token');
        showAuth();
        closeWebsocket()
    } catch (error) {
        showErrorPage(error);
    }
};