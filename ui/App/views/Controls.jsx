import React, {useEffect, useState} from "react";
import Panel from "../elements/Panel";
import Button from "../elements/Button";
import server from "../../api/resources/server";

const Controls = ({serverStatus}) => {

    const [factorioVersion, setFactorioVersion] = useState('unknown');
    const isRunning = () => {
        return serverStatus.data.status === 'running';
    }

    const startServer = () => {

    }

    const stopServer = () => {

    }

    useEffect(() => {
        (async () => {
            const res = await server.factorioVersion();
            if (res.success) {
                setFactorioVersion(res.data.version)
            }
        })();
    })

    return (
        <Panel
            title="Server Status"
            content={
                <div className="flex" slot="content">
                    <table className="w-full">
                        <thead>
                        <tr className="text-left py-1">
                            <th>Status</th>
                            <th>Address</th>
                            <th>Factorio Version</th>
                            <th>Save File</th>
                        </tr>
                        </thead>
                        <tbody>
                        <tr className="py-1">
                            <td className="pr-4">{serverStatus.data.status}</td>
                            <td className="pr-4">{serverStatus.data.address}:{serverStatus.data.port}</td>
                            <td className="pr-4">{factorioVersion}</td>
                            <td className="pr-4">{serverStatus.data.savefile}</td>
                        </tr>
                        </tbody>
                    </table>
                </div>
            }
            actions={
                <div className="flex">
                    { isRunning()
                            ? <>
                            <Button onClick={stopServer} className="mr-2" type="default">Save & Stop Server</Button>
                            <Button onClick={stopServer} type="danger">Stop Server</Button>
                        </>
                        : <Button onClick={startServer} type="success">Start Server</Button>
                    }
                </div>
            }
        />
    )
}

export default Controls;