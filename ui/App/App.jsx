import React, {useCallback, useState} from 'react';

import user from "../api/resources/user";
import Login from "./views/Login";
import {Redirect, Route} from "react-router";
import Dashboard from "./views/Dashboard";
import {BrowserRouter} from "react-router-dom";

const App = () => {

    const [isAuthenticated, setIsAuthenticated] = useState(false);

    const handleAuthenticationStatus = async () => {
        const status = await user.status();
        setIsAuthenticated(status.success);
    };

    const ProtectedRoute = useCallback(({component: Component, ...rest}) => (
        <Route {...rest} render={(props) => (
            isAuthenticated
                ? <Component {...props} />
                : <Redirect to={{
                    pathname: '/login',
                    state: {from: props.location}
                }}/>
        )}/>
    ), [isAuthenticated]);

    return (
        <BrowserRouter>
            <ProtectedRoute exact path="/" component={Dashboard}/>
            <Route path="/login" render={() => (<Login handleLogin={handleAuthenticationStatus} />)} />
        </BrowserRouter>
    );
}

export default App;
