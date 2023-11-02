import React, {useCallback, useState} from 'react';

import user from "../api/resources/user";
import Login from "./views/Login";
import {Navigate, Route, Routes} from "react-router";
import Controls from "./views/Controls";
import {BrowserRouter, Outlet} from "react-router-dom";
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

    const ProtectedRoute = ({isAuthenticated}) => {
        if (!isAuthenticated) {
            return <Navigate to="/login" state={{from: window.location.pathname}} />;
        }
        return <Outlet/>;
    }

    return (
        <BrowserRouter>
            <Routes>
                <Route path="login" element={<Login handleLogin={handleAuthenticationStatus}/>}/>

                {/* route with only `element` will cause the proper children to be place in `<Outlet/>` */}
                <Route element={<ProtectedRoute isAuthenticated={isAuthenticated}/> }>
                    <Route element={<Layout handleLogout={handleLogout} serverStatus={serverStatus} />}>
                        <Route index element={<Controls serverStatus={serverStatus}/>}/>
                        <Route path="saves" element={<Saves serverStatus={serverStatus}/>}/>
                        <Route path="mods" element={<Mods serverStatus={serverStatus}/>}/>
                        <Route path="server-settings" element={<ServerSettings serverStatus={serverStatus}/>}/>
                        <Route path="game-settings" element={<GameSettings serverStatus={serverStatus}/>}/>
                        <Route path="console" element={<Console serverStatus={serverStatus}/>}/>
                        <Route path="logs" element={<Logs serverStatus={serverStatus}/>}/>
                        <Route path="user-management" element={<UserManagement serverStatus={serverStatus}/>}/>
                        <Route path="help" element={<Help serverStatus={serverStatus}/>}/>
                    </Route>
                </Route>
            </Routes>
        </BrowserRouter>
    );
}

export default App;
