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
    },
    upload: async file => {
        let formData = new FormData();
        formData.append("savefile", file);

        const response = await client.post(`/api/saves/upload`, formData, {
            headers: {
                "Content-Type": "multipart/form-data"
            }
        });
        return response.data;
    },
    mods: async save => {
        const response = await client.post("/api/saves/mods", {
            saveFile: save
        });
        return response.data;
    }
}