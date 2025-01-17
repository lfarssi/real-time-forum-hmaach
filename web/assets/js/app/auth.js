import { registerUser, loginUser, logoutUser } from './api.js'
import { writeMessage } from './utils.js';
import { showFeed } from './feed.js';

export const showAuth = () => {
    document.body.innerHTML = ``
    const formContainer = document.createElement('div')
    formContainer.id = 'form-container'
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
                <p id="alertMsg" ></p>
                <input type="submit" value="Register">
            </form>
        </div>

        <div class="login">
            <h2>Login</h2>
            <form id="login-form" name="formLogin" action="/auth/login" method="post">
                <input type="text" name="identifier" placeholder="Email or nickname" required>
                <input type="password" name="password" placeholder="Password" required>
                <p id="login-error" ></p>
                <input type="submit" value="Login">
            </form>
        </div>
    </div>
    `
    document.body.appendChild(formContainer)
    toggleForm()
    FormSubmission()
};

const toggleForm = () => {
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
}

const FormSubmission = () => {
    const registerForm = document.querySelector('#register-form');
    const loginForm = document.querySelector('#login-form');

    [registerForm, loginForm].forEach(form => {
        form.addEventListener('submit', async (e) => {
            e.preventDefault();

            // Convert FormData to an object
            const formData = new FormData(form);
            const formObject = {};
            formData.forEach((value, key) => {
                if (key === 'age') {
                    formObject[key] = Number(value);
                } else {
                    formObject[key] = value;
                }
            });

            try {
                const response = await loginUser(formObject);
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
    });
}

const internalServerError = () => {
}

const togglePassword = () => {
    fieldPw = document.querySelector("input[name='password']");
    eye = document.querySelector('#registerForm i');
    if (fieldPw.type === 'password') {
        fieldPw.type = 'text';
        eye.classList = ['eye-off-icon'];
    } else {
        fieldPw.type = 'password';
        eye.classList = ['eye-icon'];
    }
}

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
                "password_confirmation": document.getElementById("register-password-confirmation").value,
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
                console.log(response.message);

                localStorage.removeItem("user");
                localStorage.removeItem("token");
                showAuth();

            } catch (error) {
                console.error(error);
            }
        });
    }
};