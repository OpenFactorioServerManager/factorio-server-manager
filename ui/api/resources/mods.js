import client from "../client";

const mods = {
    installed: async () => {
        const response = await client.get('/api/mods/list');
        return response.data;
    },
    toggle: async name => {
        const response = await client.post('/api/mods/toggle', {name});
        return response.data;
    },
    delete: async name => {
        const response = await client.post('/api/mods/delete', {name});
        return response.data;
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
        },
        installMultiple: async mods => {
            const response = await client.post('/api/mods/portal/install/multiple', mods);
            return response.data
        },
        install: async (downloadUrl, fileName, modName) => {
            const response = await client.post('/api/mods/portal/install', {
                downloadUrl,
                fileName,
                modName
            });
            return response.data
        },
        list: async () => {
            const response = await client.get('/api/mods/portal/list');
            return response.data
        },
        info: async mod => {
            const response = await client.get(`/api/mods/portal/info/${mod}`);
            return response.data;
        }
    }
}

export default mods;