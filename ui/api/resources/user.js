import regeneratorRuntime from "regenerator-runtime";
import client from "../client";

export default {
    status: async () => {
        const response = await client.get('/api/user/status');
        return response.data;
    },
    login: async data => {
        const response = await client.post('/api/login', data);
        return response.data;
    },
    logout: async () => {
        const response = await client.get('/api/logout');
        return response.data;
    },
    list: async () => {
        const response = await client.get('/api/user/list');
        return response.data;
    },
    add: async (user) => {
        const response = await client.post('/api/user/add', user);
        return response.data;
    },
    delete: async (username) => {
        const response = await client.post('/api/user/remove', JSON.stringify({username}));
        return response.data;
    },
    changePassword: async data => {
        const response = await client.post('/api/user/password', data);
        return response.data;
    }
}