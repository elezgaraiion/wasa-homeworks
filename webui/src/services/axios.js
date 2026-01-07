import axios from "axios";

const instance = axios.create({
    baseURL: __API_URL__,
    timeout: 5000, 
});

instance.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('userId');

        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`;
        }

        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

export default instance;