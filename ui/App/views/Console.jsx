import Panel from "../components/Panel";
import React, {useEffect, useRef, useState} from "react";
import socket from "../../api/socket";
import { useTranslation } from "react-i18next";

const Console = ({serverStatus}) => {

    const { t, i18n } = useTranslation();

    const [logs, setLogs] = useState([]);
    const consoleInput = useRef(null);

    useEffect(() => {

        const appendLog = line => {
            setLogs(lines => [...lines, line]);
        }

        socket.on('gamelog', appendLog)
        socket.emit('log subscribe')
        consoleInput.current?.focus();

        return () => {
            socket.off('gamelog', appendLog);
            socket.emit("log unsubscribe")
        }
    }, []);

    return (
        <Panel
            title={t("console.title")}
            content={
                serverStatus.running
                    ? <>
                        <ul>
                            {logs?.map((log, i) => (<li key={i}>{log}</li>))}
                        </ul>
                        <input type="text"
                               className="shadow appearance-none border w-full py-2 px-3 text-black"
                               ref={consoleInput}
                               onKeyDown={e => {
                                   if (e.key === "Enter" && socket) {
                                       socket.emit("command send", consoleInput.current.value);
                                       consoleInput.current.value = ""
                                   }
                               }}
                        />
                    </>
                    : <p className="text-red-light pt-4">
                        {t("console.error")}
                    </p>
            }
        />
    )
}

export default Console;
