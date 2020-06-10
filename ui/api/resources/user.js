import client from "../client";

export default {
    status: async () => {
        const response = await client.get('/api/user/status');
        return response.data;
    },
    login: async data => {
        const response = await client.post('/api/login');
        return response.data;
    }
}