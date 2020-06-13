import React, {useEffect, useState} from "react";
import Layout from "../components/Layout";
import log from "../../api/resources/log";

const Logs = () => {

    const [logs, setLogs] = useState([])

    useEffect(() => {
        (async () => {
            const logs = await log.tail()
            setLogs(logs.data);
        })();
    }, [])

    return (
        <div className="rounded-sm bg-gray-dark shadow-xl pb-4">
            <div className="px-4 py-2 text-xl text-dirty-white font-bold">
                Logs
            </div>
            <div className="text-white rounded-sm bg-gray-medium shadow-inner mx-4 px-6 pt-4 pb-6">
                <ul>
                    {logs.map(log => (<li key={log}>{log}</li>))}
                </ul>
            </div>
        </div>
    );
}

export default Logs;