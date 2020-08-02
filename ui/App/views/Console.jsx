import Panel from "../components/Panel";
import React, {useEffect, useRef, useState} from "react";
import Socket from "../../api/socket"

const Console = () => {

    const socket = new Socket();

    const [logs, setLogs] = useState([]);
    const consoleInput = useRef(null);

    useEffect(() => {
        // todo: maybe check for ready state
        socket.emit("log subscribe")
        socket.on('log update', log => {
            console.log(log)
            setLogs((logs) => {
                [...logs].push(log)
            });
        })
        consoleInput.current?.focus();
    }, []);

    return (
        <Panel
            title="Console"
            content={
                <>
                    <ul>
                        {logs?.map(log => (<li key={log}>{log}</li>))}
                    </ul>
                    <input type="text" ref={consoleInput} className="shadow appearance-none border w-full py-2 px-3 text-black" onKeyPress={e => {
                            if (e.key === "Enter") {
                                socket.emit("command send", consoleInput.current.value);
                                consoleInput.current.value = ""
                            }
                        }
                    }/>
                </>
            }
        />
    )
}

export default Console;