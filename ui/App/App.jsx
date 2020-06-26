import React, {useCallback, useState} from 'react';

import user from "../api/resources/user";
import Login from "./views/Login";
import {Redirect, Route, Switch, useHistory} from "react-router";
import Controls from "./views/Controls";
import {BrowserRouter} from "react-router-dom";
import Logs from "./views/Logs";
import Saves from "./views/Saves";
import Layout from "./components/Layout";
import server from "../api/resources/server";
import Mods from "./views/Mods";
import UserManagement from "./views/UserManagment";

const App = () => {

    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [serverStatus, setServerStatus] = useState(null);
    const history = useHistory();

    const updateServerStatus = async () => {
        const status = await server.status();
        if (status.success) {
            setServerStatus(status)
        }
    }

    const handleAuthenticationStatus = async () => {
        const status = await user.status();
        setIsAuthenticated(status.success);
        await updateServerStatus();
    };

    const handleLogout = async () => {
        const loggedOut = await user.logout();
        if (loggedOut.success) {
            setIsAuthenticated(false);
            history.push('/login');
        }
    }

    const ProtectedRoute = useCallback(({component: Component, ...rest}) => (
        <Route {...rest} render={(props) => (
            isAuthenticated
                ? <Component serverStatus={serverStatus} updateServerStatus={updateServerStatus} {...props} />
                : <Redirect to={{
                    pathname: '/login',
                    state: {from: props.location}
                }}/>
        )}/>
    ), [isAuthenticated]);

    return (
        <BrowserRouter>
            <Switch>
                <Route path="/login" render={() => (<Login handleLogin={handleAuthenticationStatus}/>)}/>

                <Layout handleLogout={handleLogout} serverStatus={serverStatus} updateServerStatus={updateServerStatus}>
                    <ProtectedRoute exact path="/" component={Controls}/>
                    <ProtectedRoute path="/saves" component={Saves}/>
                    <ProtectedRoute path="/mods" component={Mods}/>
                    <ProtectedRoute path="/server-settings" component={Controls}/>
                    <ProtectedRoute path="/game-settings" component={Controls}/>
                    <ProtectedRoute path="/console" component={Controls}/>
                    <ProtectedRoute path="/logs" component={Logs}/>
                    <ProtectedRoute path="/user-management" component={UserManagement}/>
                    <ProtectedRoute path="/help" component={Controls}/>
                </Layout>
            </Switch>
        </BrowserRouter>
    );
}

export default App;
