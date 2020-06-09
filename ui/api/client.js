import Axios from "axios";

const client = Axios.create({
    withCredentials: true
});

export default client;