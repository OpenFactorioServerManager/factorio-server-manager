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
    update: async ({modName, downloadUrl, fileName}) => {
        const response = await client.post('/api/mods/update', {modName, downloadUrl, fileName})
        return response.data;
    },
    upload: async file => {
        let formData = new FormData();
        formData.append("mod_file", file);

        const response = await client.post('/api/mods/upload', formData, {
            headers: {
                "Content-Type": "multipart/form-data"
            }
        });
        return response.data;
    },
    deleteAll: async () => {
        const response = await client.post('/api/mods/delete/all');
        return response.data;
    },
    downloadAllURL: '/api/mods/download',
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
    },
    packs: {
        list: async () => {
            const response = await client.get('/api/mods/packs/list');
            return response.data;
        },
        create: async name => {
            const response = await client.post('/api/mods/packs/create', {name});
            return response.data;
        },
        delete: async name => {
            const response = await client.post(`/api/mods/packs/${name}/delete`);
            return response.data;
        },
        download: async name => {
            const response = await client.get(`/api/mods/packs/${name}/download`);
            return response.data;
        },
        load: async name => {
            const response = await client.post(`/api/mods/packs/${name}/load`);
            return response.data;
        },
        mods: {
            list: async packName => {
                const response = await client.get(`/api/mods/packs/${packName}/list`);
                return response.data;
            },
            toggle: async (packName, modName) => {
                const response = await client.post(`/api/mods/packs/${packName}/mod/toggle`, {
                    name: modName
                });
                return response.data;
            },
            update: async (packName, {modName, downloadUrl, fileName}) => {
                const response = await client.post(`/api/mods/packs/${packName}/mod/update`, {modName, downloadUrl, fileName})
                return response.data;
            },
            delete: async (packName, modName) => {
                const response = await client.post(`/api/mods/packs/${packName}/mod/delete`, {
                    name: modName
                });
                return response.data;
            },
        }
    }
}

export default mods;