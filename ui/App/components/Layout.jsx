import React, {useEffect, useState} from "react";
import server from "../../api/resources/server";
import {NavLink, useHistory} from "react-router-dom";
import Button from "../elements/Button";

const Layout = ({children, handleLogout}) => {

    const [serverStatus, setServerStatus] = useState(null);
    const history = useHistory();

    useEffect(() => {
        (async () => {
            const status = await server.status();
            if (status.success) {
                setServerStatus(status)
            }
        })();
    }, []);

    const Status = ({info}) => {

        let text = 'Unknown';
        let color = 'gray-light';

        if (info && info.success) {
            if (info.data.status === 'running') {
                text = 'Running';
                color = 'green';
            } else if (info.data.status === 'stopped') {
                text = 'Stopped';
                color = 'red';
            }
        }

        return (
            <div className={`bg-${color} rounded-sm px-2 py-1 text-black`}>{text}</div>
        )
    }

    const Link = ({children, to, last}) => {
        return (
            <NavLink
                exact={true}
                to={to}
                activeClassName="bg-orange"
                className={`bg-gray-light hover:bg-orange text-black font-bold py-2 px-4 w-full block${last ? '' : ' mb-1'}`}
            >{children}</NavLink>)
    }

    return (
        <div className="flex md:flex-row-reverse flex-wrap">

            {/*Main*/}
            <div className="w-full md:w-5/6 bg-gray-100 bg-banner bg-fixed min-h-screen">
                <div className="container mx-auto bg-gray-100 pt-16 px-6">
                    {children}
                </div>
            </div>

            {/*Sidebar*/}
            <div
                className="w-full border-r border-black md:w-1/6 bg-gray-dark fixed bottom-0 md:top-0 md:left-0 h-16 md:h-screen">
                <div className="py-4 px-2 border-b-2 border-black items-center text-center">
                    <img src="/images/factorio.jpg" className="inline h-8" alt="Factorio Logo"/>
                    <span className="text-dirty-white pl-2 text-xl">Factorio Server Manager</span>
                </div>
                <div className="py-4 px-2 border-b-2 border-black">
                    <h1 className="text-dirty-white text-lg mb-2 mx-4">Server Status</h1>
                    <div className="mx-4 mb-4 text-center">
                        <Status info={serverStatus}/>
                    </div>
                </div>
                <div className="py-4 px-2 border-b-2 border-black">
                    <h1 className="text-dirty-white text-lg mb-2 mx-4">Server Management</h1>
                    <div className="text-white text-center rounded-sm bg-black shadow-inner mx-4 p-1">
                        <Link to="/">Controls</Link>
                        <Link to="/saves">Saves</Link>
                        <Link to="/mods">Mods</Link>
                        <Link to="/server-settings">Server Settings</Link>
                        <Link to="/game-settings">Game Settings</Link>
                        <Link to="/console">Console</Link>
                        <Link to="/logs" last={true}>Logs</Link>
                    </div>
                </div>
                <div className="py-4 px-2 border-b-2 border-black">
                    <h1 className="text-dirty-white text-lg mb-2 mx-4">Administration</h1>
                    <div className="text-white text-center rounded-sm bg-black shadow-inner mx-4 p-1">
                        <Link to="/user-management">User Management</Link>
                        <Link to="/help" last={true}>Help</Link>
                    </div>
                </div>
                <div className="py-4 px-2 border-b-2 border-black">
                    <div className="text-white text-center rounded-sm bg-black shadow-inner mx-4 p-1">
                        <Button type="danger" onClick={handleLogout}>Logout</Button>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default Layout;