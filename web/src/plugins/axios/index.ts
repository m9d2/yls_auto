import axios from 'axios';

let baseUrl = import.meta.env.MODE == "release" ? window.location.origin : import.meta.env.VITE_SERVER
baseUrl = baseUrl + "/v1"

const request = axios.create({
    baseURL: baseUrl,
    timeout: 10000,
    responseType: "json",
    headers: {
        'Content-Type': 'application/json',
    }
});

let loadingInstance: any
request.interceptors.request.use(
    config => {
        return config;
    },
    error => {
        return Promise.reject(error);
    }
);

request.interceptors.response.use(
    response => {
        const { code, message } = response.data;
        if (loadingInstance) {
            loadingInstance.close();
        }
        if (code === 200) {
            return response.data;
        } else if (code == 400) {
            throw new Error(message)
        }
    },
    error => {
        return Promise.reject(error);
    },

);

export default request;
