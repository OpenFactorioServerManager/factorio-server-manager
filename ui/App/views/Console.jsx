import Panel from "../components/Panel";
import React, {useEffect, useRef, useState} from "react";
import socket from "../../api/socket";
import Input from "../components/Input";

const Console = ({serverStatus}) => {

    const [logs, setLogs] = useState([]);
    const consoleInput = useRef(null);
    const isRunning = serverStatus.status === 'running';

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
            title="Console"
            content={
                isRunning
                    ? <>
                        <ul>
                            {logs?.map((log, i) => (<li key={i}>{log}</li>))}
                        </ul>
                        <Input type="text"
                               inputRef={consoleInput}
                               onKeyDown={e => {
                                   if (e.key === "Enter" && socket) {
                                       socket.emit("command send", consoleInput.current.value);
                                       consoleInput.current.value = ""
                                   }
                               }}
                        />
                    </>
                    : <p className="text-red-light pt-4">
                        The console is not available, because Factorio is not running.
                    </p>
            }
        />
    )
}

export default Console;