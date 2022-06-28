import React, {useCallback, useState} from 'react';

import user from "../api/resources/user";
import Login from "./views/Login";
import {Navigate, Routes} from "react-router";
import Controls from "./views/Controls";
import {BrowserRouter, Link} from "react-router-dom";
import Logs from "./views/Logs";
import Saves from "./views/Saves/Saves";
import Layout from "./components/Layout";
import server from "../api/resources/server";
import Mods from "./views/Mods/Mods";
import UserManagement from "./views/UserManagement/UserManagment";
import ServerSettings from "./views/ServerSettings";
import GameSettings from "./views/GameSettings";
import Console from "./views/Console";
import Help from "./views/Help";
import socket from "../api/socket";
import {Flash} from "./components/Flash";


const App = () => {

    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [serverStatus, setServerStatus] = useState(null);

    const handleAuthenticationStatus = useCallback(async (status) => {
        if (status?.username) {
            setIsAuthenticated(true);

            const status = await server.status();
            setServerStatus(status);

            socket.emit('server status subscribe');
            socket.on('server_status', status => {
                setServerStatus(JSON.parse(status));
            });
        }
    },[]);

    const handleLogout = useCallback(async () => {
        const loggedOut = await user.logout();
        if (loggedOut) {
            setIsAuthenticated(false);
        }
    }, []);

    const ProtectedRoute = useCallback(({component: Component, ...rest}) => (
        <Link {...rest} render={(props) => (
            isAuthenticated && Component
                ? <Component serverStatus={serverStatus} {...props} />
                : <Navigate to={{
                    pathname: '/login',
                    state: {from: props.location}
                }}/>
        )}/>
    ), [isAuthenticated, serverStatus]);

    return (
        <BrowserRouter basename="/">
            <Routes>
                <Link to="/login" render={() => (<Login handleLogin={handleAuthenticationStatus}/>)}/>

                <Layout handleLogout={handleLogout} serverStatus={serverStatus}>
                    <ProtectedRoute end to="/" component={Controls}/>
                    <ProtectedRoute to="/saves" component={Saves}/>
                    <ProtectedRoute to="/mods" component={Mods}/>
                    <ProtectedRoute to="/server-settings" component={ServerSettings}/>
                    <ProtectedRoute to="/game-settings" component={GameSettings}/>
                    <ProtectedRoute to="/console" component={Console}/>
                    <ProtectedRoute to="/logs" component={Logs}/>
                    <ProtectedRoute to="/user-management" component={UserManagement}/>
                    <ProtectedRoute to="/help" component={Help}/>
                    <Flash/>
                </Layout>
            </Routes>
        </BrowserRouter>
    );
}

export default App;
