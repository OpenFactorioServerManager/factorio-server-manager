import React, {useEffect, useState} from "react";
import Panel from "../components/Panel";
import Button from "../components/Button";
import server from "../../api/resources/server";
import savesResource from "../../api/resources/saves";
import {useForm} from "react-hook-form";

const Controls = ({serverStatus, updateServerStatus}) => {

    const [factorioVersion, setFactorioVersion] = useState('unknown');
    const isRunning = serverStatus.data.status === 'running';
    const [saves, setSaves] = useState([]);

    const { handleSubmit, register, errors } = useForm();

    const startServer = async (data) => {
        await server.start(data.ip, parseInt(data.port), data.save);
        await updateServerStatus();
    }

    const stopServer = async () => {
        await server.stop();
        await updateServerStatus();
    }

    const killServer = async () => {
        await server.kill();
        await updateServerStatus();
    }

    useEffect(() => {
        server.factorioVersion()
            .then(res => {
                if (res.success) {
                    setFactorioVersion(res.data.version)
                }
            });

        savesResource.list()
            .then(res => {
                if (res.success) {
                    setSaves(res.data)
                }
            });
    }, [])

    return (
        <form onSubmit={handleSubmit(startServer)}>
        <Panel
            title="Server Status"
            content={
                <div className="flex">
                    <table className="w-full">
                        <thead>
                        <tr className="text-left py-1">
                            <th>Status</th>
                            <th>IP</th>
                            <th>Port</th>
                            <th>Factorio Version</th>
                            <th>Save File</th>
                        </tr>
                        </thead>
                        <tbody>
                        {isRunning
                            ? <tr className="py-1">
                                <td className="pr-4 py-2">{serverStatus.data.status}</td>
                                <td className="pr-4 py-2">{serverStatus.data.address}</td>
                                <td className="pr-4 py-2">{serverStatus.data.port}</td>
                                <td className="pr-4 py-2">{factorioVersion}</td>
                                <td className="pr-4 py-2">{serverStatus.data.savefile}</td>
                            </tr>
                            : <tr className="py-1">
                                <td className="pr-4 py-2">{serverStatus.data.status}</td>
                                <td className="pr-4">
                                    <input
                                        name="ip"
                                        className="shadow appearance-none w-full py-2 px-3 text-black"
                                        type="text"
                                        defaultValue={"0.0.0.0"}
                                        ref={register({required: true, pattern: '^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$'})}
                                    />
                                    {errors.ip && <span className="block text-red">IP is required and must be valid.</span>}
                                </td>
                                <td className="pr-4">
                                    <input
                                        name="port"
                                        className="shadow appearance-none w-full py-2 px-3 text-black"
                                        type="number"
                                        min={1}
                                        defaultValue={"34197"}
                                        ref={register({required: true})}
                                    />
                                    {errors.port && <span className="block text-red">Port is required</span>}
                                </td>
                                <td className="pr-4 py-2">{factorioVersion}</td>
                                <td className="pr-4 py-2">
                                    <div className="relative">
                                        <select
                                            name="save"
                                            className="block shadow appearance-none w-full py-2 px-3 text-black"
                                            ref={register({required: true})}
                                        >
                                            {saves.map(save => (
                                                <option value={save.name} key={save.name}>{save.name}</option>))}
                                        </select>
                                        <div
                                            className="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-black">
                                            <svg className="fill-current h-4 w-4" xmlns="http://www.w3.org/2000/svg"
                                                 viewBox="0 0 20 20">
                                                <path
                                                    d="M9.293 12.95l.707.707L15.657 8l-1.414-1.414L10 10.828 5.757 6.586 4.343 8z"/>
                                            </svg>
                                        </div>
                                    </div>
                                </td>
                            </tr>
                        }
                        </tbody>
                    </table>
                </div>
            }
            actions={
                <div className="flex">
                    {isRunning
                        ? <>
                            <Button onClick={stopServer} className="mr-2" type="default">Save & Stop Server</Button>
                            <Button onClick={killServer} type="danger">Kill Server</Button>
                        </>
                        : <Button isSubmit={true} type="success">Start Server</Button>
                    }
                </div>
            }
        />
        </form>
    )
}

export default Controls;