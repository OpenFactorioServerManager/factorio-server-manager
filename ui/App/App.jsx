import React, {useCallback, useEffect, useState} from 'react';

import user from "../api/resources/user";
import Login from "./views/Login";
import {Redirect, Route, Switch} from "react-router";
import Dashboard from "./views/Dashboard";
import {BrowserRouter} from "react-router-dom";

const App = () => {

    const [isAuthenticated, setIsAuthenticated] = useState(false);

    const handleLogin = async () => {
        const status = await user.status();
        setIsAuthenticated(status.success);
    };

    // on mount check if user is authenticated
    useEffect(() => {

    }, [])



    const ProtectedRoute = useCallback(({component: Component, loggedIn, ...rest}) => (
        <Route {...rest} render={(props) => (
            loggedIn
                ? <Component {...props} />
                : <Redirect to={{
                    pathname: '/login',
                    state: {from: props.location}
                }}/>
        )}/>
    ), [isAuthenticated]);

    // List of all Routes
    return (
        <BrowserRouter>
            <ProtectedRoute exact path="/" loggedIn={isAuthenticated} component={Dashboard}/>
            <Route path="/login" component={Login} handleLogin={handleLogin()}/>
        </BrowserRouter>
    );
}

export default App;
