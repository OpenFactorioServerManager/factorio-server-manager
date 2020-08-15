import React, {useEffect, useState} from "react";
import saveClient from "../../../api/resources/saves";
import Panel from "../../components/Panel";
import ButtonLink from "../../components/ButtonLink";
import Button from "../../components/Button";
import CreateSaveForm from "./components/CreateSaveForm";
import UploadSaveForm from "./components/UploadSaveForm";

const Saves = ({serverStatus}) => {

    const [saves, setSaves] = useState([]);

    const updateList = async () => {
        const res = await saveClient.list();
        if (res.success) {
            setSaves(res.data);
        }
    }

    useEffect(() => {
        updateList()
    }, []);

    const deleteSave = async (save) => {
        const res = await saveClient.delete(save);
        if (res.success) {
            updateList()
        }
    }

    return (
        <>
            <div className="flex mb-6">
                <Panel
                    title="Create Save"
                    className="w-1/2 mr-3"
                    content={
                        serverStatus.data.status === "running"
                            ? <p className="text-red-light pt-4 pb-24">
                                Create a new Save is only possible if the Factorio server is
                                not running.
                            </p>
                            : <CreateSaveForm onSuccess={updateList}/>
                    }
                />
                <Panel
                    title="Upload Save"
                    className="w-1/2 ml-3"
                    content={<UploadSaveForm onSuccess={updateList}/>}
                />
            </div>

            <Panel
                className="mb-4"
                title="Saves"
                content={
                    <table className="w-full">
                        <thead>
                        <tr className="text-left py-1">
                            <th>Name</th>
                            <th>Last Modified At</th>
                            <th>Size</th>
                            <th>Actions</th>
                        </tr>
                        </thead>
                        <tbody>
                        {saves.map(save =>
                            <tr className="py-1" key={save.name}>
                                <td className="pr-4">{save.name}</td>
                                <td className="pr-4">{(new Date(save.last_mod)).toISOString().replace('T', ' ').split('.')[0]}</td>
                                <td>{parseFloat(save.size / 1024 / 1024).toFixed(3)} MB</td>
                                <td>
                                    <ButtonLink size="sm" href={`/api/saves/dl/${save.name}`}
                                                className="mr-2">Download</ButtonLink>
                                    <Button size="sm" onClick={() => deleteSave(save)} type="danger">Delete</Button>
                                </td>
                            </tr>
                        )}
                        </tbody>
                    </table>
                }
            />
        </>
    )
}

export default Saves;