import client from "../client";

const mods = {
    installed: async () => {
        const response = await client.get('/api/mods/list/installed');
        return response.data;
    },
    toggle: async modName => {
        let data = new FormData();
        data.set('modName', modName);

        const response = await client.post('/api/mods/toggle', data);
        return response.data;
    },
    delete: async modName => {
        const data = new FormData();
        data.set('modName', modName);

        const response = await client.post('/api/mods/delete', data);
        return response.data;
    },
    details: async modName => {
        const data = new FormData();
        data.set('modId', modName);

        const response = await client.post('/api/mods/details', data);

        return {
            success: response.data.success,
            data: JSON.parse(response.data.data)
        };
    },
    update: async (modName, downloadUrl, fileName) => {
        const data = new FormData();
        data.set('modName', modName);
        data.set('downloadUrl', downloadUrl);
        data.set('filename', fileName);

        const response = await client.post('/api/mods/update', data, {
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            }
        })
        return response.data;
    },
    deleteAll: async () => {
        const response = await client.post('/api/mods/delete/all');
        return response.data;
    }
}

export default mods;