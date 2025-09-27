import React, {useEffect, useState} from "react";
import log from "../../api/resources/log";
import Panel from "../components/Panel";
import { useTranslation } from "react-i18next";

const Logs = () => {

    const { t, i18n } = useTranslation();
    
    const [logs, setLogs] = useState([])

    useEffect(() => {
        (async () => {
            const logs = await log.tail()
            setLogs(logs);
        })();
    }, [])

    return (
        <Panel
            title={t("logs.title")}
            content={
                <ul>
                    {logs.map((log,index) => (<li key={index}>{log}</li>))}
                </ul>
            }
        />
    );
}

export default Logs;