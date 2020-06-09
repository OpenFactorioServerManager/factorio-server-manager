import React, {useEffect, useState} from 'react';

import user from "../api/resources/user";
import Login from "./views/Login";
import UsersContent from "./components/UsersContent";

const App = () => {

    const [isAuthenticated, setIsAuthenticated] = useState(false);

    // on mount check if user is authenticated
    useEffect( () => {
        (async () => {
            const status = await user.status();
            setIsAuthenticated(status.success);
        })();
    }, [])


    return isAuthenticated ? (<UsersContent/>) : (<Login/>);
}

export default App;
