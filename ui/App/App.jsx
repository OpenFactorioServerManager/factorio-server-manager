import React, {useCallback, useState} from 'react';

import user from "../api/resources/user";
import Login from "./views/Login";
import {Redirect, Route} from "react-router";
import Controls from "./views/Controls";
import {BrowserRouter} from "react-router-dom";
import Logs from "./views/Logs";
import Saves from "./views/Saves";

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
            <ProtectedRoute exact path="/" component={Controls}/>
            <ProtectedRoute path="/saves" component={Saves}/>
            <ProtectedRoute path="/mods" component={Controls}/>
            <ProtectedRoute path="/server-settings" component={Controls}/>
            <ProtectedRoute path="/game-settings" component={Controls}/>
            <ProtectedRoute path="/console" component={Controls}/>
            <ProtectedRoute path="/logs" component={Logs}/>
            <ProtectedRoute path="/user-management" component={Controls}/>
            <Route path="/login" render={() => (<Login handleLogin={handleAuthenticationStatus}/>)}/>
        </BrowserRouter>
    );
}

export default App;
