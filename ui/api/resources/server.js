import client from "../client";

export default {
    factorioVersion: async () => {
        const response = await client.get('/api/server/facVersion');
        return response.data;
    },
    status: async () => {
        const response = await client.get('/api/server/status');
        return response.data;
    },
}