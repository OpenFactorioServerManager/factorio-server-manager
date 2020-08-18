import Panel from "../components/Panel";
import React, {useEffect, useRef, useState} from "react";

const Console = ({serverStatus}) => {

    const [logs, setLogs] = useState([]);
    const consoleInput = useRef(null);
    const [socket, setSocket] = useState(null);
    const isRunning = serverStatus.data.status === 'running';

    useEffect(() => {
        let ws_scheme = window.location.protocol === "https:" ? "wss" : "ws";
        const socket = new WebSocket(ws_scheme + "://" + window.location.host + "/ws");

        socket.onopen = function (e) {
            console.log("Socket Open")
            socket.send(JSON.stringify({name: "log subscribe"}));
            setSocket(socket)
        };

        socket.onmessage = e => {
            const {name, data} = JSON.parse(e.data)
            console.log(name)
            switch (name) {
                case 'log update':
                    setLogs(oldLogs => [...oldLogs, data])
                    break;
                default:
                    break;
            }
        }

        socket.onerror = e => {
            console.log(e)
        }

        consoleInput.current?.focus();

        return () => {
            socket.close();
        }

    }, []);

    return (
        <Panel
            title="Console"
            content={
                <>
                    {isRunning
                    ? <>
                        <ul>
                            {logs?.map((log, i) => (<li key={i}>{log}</li>))}
                        </ul>
                        <input type="text" ref={consoleInput}
                               className="shadow appearance-none border w-full py-2 px-3 text-black" onKeyPress={e => {
                            if (e.key === "Enter" && socket) {
                                const message = {name: "command send", data: consoleInput.current.value};
                                socket.send(JSON.stringify(message));
                                consoleInput.current.value = ""
                            }
                        }
                        }/>
                    </>
                        : <p className="text-red-light pt-4">
                            The console is not available, because Factorio is not running.
                        </p>
                    }

                </>
            }
        />
    )
}

export default Console;