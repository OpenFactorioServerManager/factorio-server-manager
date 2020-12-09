import React, {useEffect, useState} from "react";
import log from "../../api/resources/log";
import Panel from "../components/Panel";

const Logs = () => {

    const [logs, setLogs] = useState([])

    useEffect(() => {
        (async () => {
            const logs = await log.tail()
            setLogs(logs);
        })();
    }, [])

    return (
        <Panel
            title="Logs"
            content={
                <ul>
                    {logs.map((log,index) => (<li key={index}>{log}</li>))}
                </ul>
            }
        />
    );
}

export default Logs;