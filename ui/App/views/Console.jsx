import Panel from "../components/Panel";
import React, {useEffect, useRef, useState} from "react";
import socket from "../../api/socket";

const Console = ({serverStatus}) => {

    const [logs, setLogs] = useState([]);
    const consoleInput = useRef(null);
    const isRunning = serverStatus.data.status === 'running';

    useEffect(() => {

        const appendLog = line => {
            setLogs(lines => [...lines, line]);
        }

        socket.on('log update', appendLog)
        socket.emit('log subscribe')
        consoleInput.current?.focus();

        return () => {
            socket.off('log update', appendLog);
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
                        <input type="text" ref={consoleInput}
                               className="shadow appearance-none border w-full py-2 px-3 text-black" onKeyPress={e => {
                            if (e.key === "Enter" && socket) {
                                socket.emit("command send", consoleInput.current.value);
                                consoleInput.current.value = ""
                            }
                        }
                        }/>
                    </>
                    : <p className="text-red-light pt-4">
                        The console is not available, because Factorio is not running.
                    </p>
            }
        />
    )
}

export default Console;