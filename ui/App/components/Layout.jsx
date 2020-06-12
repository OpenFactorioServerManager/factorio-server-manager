import React, {useEffect, useRef, useState} from "react";
import server from "../../api/resources/server";

const Layout = (props) => {

    const [serverStatus, setServerStatus] = useState(null);

    useEffect(() => {
        (async () => {
            const status = await server.status();
            if (status.success) {
                setServerStatus(status)
            }
        })();
    }, []);

    const Status = (props) => {

        let text = 'Unknown';
        let color = 'gray-light';

        if (props.info && props.info.success) {
            console.log(props.info);
            if (props.info.data.status === 'running') {
                text = 'Running';
                color = 'green';
            } else if (props.info.data.status === 'stopped') {
                text = 'Stopped';
                color = 'red';
            }
        }

        return (
            <div className={`bg-${color} rounded-sm px-2 py-1 text-black`}>{text}</div>
        )
    }

    return (
        <div className="flex md:flex-row-reverse flex-wrap">

            {/*Main*/}
            <div className="w-full md:w-4/5 bg-gray-100 bg-banner min-h-screen">
                <div className="container bg-gray-100 pt-16 px-6">
                    {props.children}
                </div>
            </div>

            {/*Sidebar*/}
            <div
                className="w-full border-r border-black md:w-1/5 bg-gray-dark fixed bottom-0 md:top-0 md:left-0 h-16 md:h-screen">
                <div className="py-4 px-2 border-b-2 border-black items-center text-center">
                    <img src="/images/factorio.jpg" className="inline h-8" alt="Factorio Logo"/>
                    <span className="text-dirty-white pl-2 text-xl">Factorio Server Manager</span>
                </div>
                <div className="py-4 px-2 border-b-2 border-black">
                    <h1 className="text-dirty-white text-lg mb-2 mx-4">Server Status</h1>
                    <div className="text-white text-center rounded-sm bg-gray-medium shadow-inner mx-4 px-6 py-6 mb-4">
                        <Status info={serverStatus}/>
                    </div>
                </div>
                <nav className="py-4 px-2 border-b-2 border-black">
                    <h1>Server Control</h1>
                </nav>
            </div>
        </div>
    );
}

export default Layout;