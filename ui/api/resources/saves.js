import client from "../client";

export default {
    list: async () => {
        const response = await client.get('/api/saves/list');
        return response.data;
    },
}