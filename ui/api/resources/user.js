import client from "../client";

export default {
    status: async () => {
        const response = await client.get('/api/user/status');
        return response.data;
    }
}