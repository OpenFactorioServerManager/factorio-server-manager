import client from "../client";

export default {
    server: {
        list: async () => {
            const response = await client.get('/api/settings')
            return response.data;
        },
        update: async data => {
            const response = await client.post('/api/settings/update', data)
            return response.data;
        }
    },
    game: {
        list: async () => {
            const response = await client.get('/api/config');
            return response.data;
        }
    }
}