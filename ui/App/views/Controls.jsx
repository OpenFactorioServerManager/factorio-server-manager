import React, {useEffect, useState} from "react";
import Panel from "../components/Panel";
import Button from "../components/Button";
import server from "../../api/resources/server";
import savesResource from "../../api/resources/saves";
import {useForm} from "react-hook-form";
import Select from "../components/Select";
import Input from "../components/Input";

const Controls = ({serverStatus, updateServerStatus}) => {

    const [factorioVersion, setFactorioVersion] = useState('unknown');
    const isRunning = serverStatus.status === 'running';
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
            .then(res => setFactorioVersion(res.version));

        savesResource.list()
            .then(res => setSaves(res));
    }, [])

    return (
        <form onSubmit={handleSubmit(startServer)}>
        <Panel
            title="Server Status"
            content={
                <div className="lg:flex">
                    { isRunning
                        ? <>
                            <div className="lg:w-1/5 mb-2">
                                <div className="font-bold">Status</div>
                                <div>{serverStatus.status}</div>
                            </div>
                            <div className="lg:w-1/5 mb-2">
                                <div className="font-bold">IP</div>
                                <div>{serverStatus.address}</div>
                            </div>
                            <div className="lg:w-1/5 mb-2">
                                <div className="font-bold">Port</div>
                                <div>{serverStatus.port}</div>
                            </div>
                            <div className="lg:w-1/5 mb-2">
                                <div className="font-bold">Factorio Version</div>
                                <div>{factorioVersion}</div>
                            </div>
                            <div className="lg:w-1/5 mb-2">
                                <div className="font-bold">Save</div>
                                <div>{serverStatus.savefile}</div>
                            </div>
                        </>
                        : <>
                            <div className="lg:w-1/5 mb-2">
                                <div className="font-bold">Status</div>
                                <div>{serverStatus.status}</div>
                            </div>
                            <div className="lg:w-1/5 mb-2 mr-0 lg:mr-4">
                                <div className="font-bold">IP</div>
                                <Input
                                    name="ip"
                                    defaultValue={"0.0.0.0"}
                                    inputRef={register({required: true, pattern: '^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$'})}
                                />
                                {errors.ip && <span className="block text-red">IP is required and must be valid.</span>}
                            </div>
                            <div className="lg:w-1/5 mb-2 mr-0 lg:mr-4">
                                <div className="font-bold">Port</div>
                                <Input
                                    name="port"
                                    type="number"
                                    min={1}
                                    defaultValue={"34197"}
                                    inputRef={register({required: true})}
                                />
                                {errors.port && <span className="block text-red">Port is required</span>}
                            </div>
                            <div className="lg:w-1/5 mb-2 mr-0 lg:mr-4">
                                <div className="font-bold">Factorio Version</div>
                                <div>{factorioVersion}</div>
                            </div>
                            <div className="lg:w-1/5 mb-2">
                                <div className="font-bold">Save</div>
                                <div className="relative">
                                    <Select
                                        name="save"
                                        inputRef={register({required: true})}
                                        defaultValue="Load Latest"
                                        options={saves.map(save => new Object({
                                            value: save.name,
                                            name: save.name
                                        }))}
                                    />
                                </div>
                            </div>
                        </>
                    }
                </div>
            }
            actions={
                <div className="md:flex">
                    {isRunning
                        ? <>
                            <Button onClick={stopServer} size="sm" className="w-full md:w-auto mb-2 md:mb-0 md:mr-2" type="default">Save & Stop Server</Button>
                            <Button onClick={killServer} size="sm" type="danger" className="w-full md:w-auto">Kill Server</Button>
                        </>
                        : <Button isSubmit={true} size="sm" type="success" className="w-full md:w-auto">Start Server</Button>
                    }
                </div>
            }
        />
        </form>
    )
};

export default Controls;