import {
    registerUser,
    loginUser,
    logoutUser
} from './api.js'

const token = localStorage.getItem('token')

export const submitRegistration = async () => {
    try {
        const user = {
            "first_name": "Hamza",
            "last_name": "Maach",
            "email": "hamazmaaaaaaaach56@gmail.com",
            "nickname": "hmaaaaaaaaaaach",
            "gender": "male",
            "age": 22,
            "password": "Hamza.1234",
            "password_confirmation": "Hamza.1234"
        }

        const message = await registerUser(user)
        console.log("registration message: ", message)
    } catch (error) {
        console.log(error)
    }
}

export const submitLogin = async (credentials) => {
    try {
        const response = await loginUser(credentials)
        localStorage.setItem("user", JSON.stringify(response.user))
        localStorage.setItem("token", response.token)
    } catch (error) {
        console.log(error)
    }
}

export const submitLogout = async () => {
    try {
        const message = await logoutUser(token)
        console.log("logout message: ", message)
        localStorage.removeItem("user");
        localStorage.removeItem("token");
    } catch (error) {
        console.log(error)
    }
}
