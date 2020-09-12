import client from "../client";

const mods = {
    installed: async () => {
        const response = await client.get('/api/mods/list');
        return response.data;
    },
    toggle: async modName => {
        const response = await client.post('/api/mods/toggle', JSON.stringify({modName}));
        return response.data;
    },
    delete: async modName => {
        const response = await client.post('/api/mods/delete', JSON.stringify({modName}));
        return response.data;
    },
    details: async modName => {
        const response = await client.post('/api/mods/details', JSON.stringify({modName}));

        return {
            data: JSON.parse(response.data)
        };
    },
    update: async (modName, downloadUrl, fileName) => {
        const data = {
            modName: modName,
            downloadUrl: downloadUrl,
            fileName: fileName,
        }

        const response = await client.post('/api/mods/update', data)
        return response.data;
    },
    deleteAll: async () => {
        const response = await client.post('/api/mods/delete/all');
        return response.data;
    },
    portal: {
        login: async (username, password) => {
            // const data = new FormData();
            // data.set('username', username);
            // data.set('password', password);

            const response = await client.post('/api/mods/portal/login', {
                username,
                password
            });
            return response.data;
        },
        status: async () => {
            const response = await client.get('/api/mods/portal/loginstatus');
            return response.data;
        },
        logout: async () => {
            const response = await client.get('/api/mods/portal/logout');
            return response.data
        }
    }
}

export default mods;