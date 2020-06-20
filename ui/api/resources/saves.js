import client from "../client";

export default {
    list: async () => {
        const response = await client.get('/api/saves/list');
        return response.data;
    },
    delete: async (save) => {
        const response = await client.get(`/api/saves/rm/${save.name}`);
        return response.data;
    },
    create: async (name) => {
        const response = await client.get(`/api/saves/create/${name}`);
        return response.data;
    }
}