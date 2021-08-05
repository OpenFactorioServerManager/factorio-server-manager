import client from "../client";

export default {
    tail: async () => {
        const response = await client.get('/api/log/tail');
        return response.data;
    },
}