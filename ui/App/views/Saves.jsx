import React, {useEffect, useState} from "react";
import save from "../../api/resources/saves";
import Panel from "../elements/Panel";

const Saves = () => {

    const [saves, setSaves] = useState([]);

    useEffect(() => {
        (async () => {
            const list = await save.list();
            if (list.success) {
                console.log(list)
                setSaves(list.data);
            }
        })()
    }, [])

    return (
        <Panel
            title="Saves"
            content={<table className="w-full">
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
                        <td className="pr-4">{save.last_mod}</td>
                        <td>{parseFloat(save.size / 1024 / 1024).toFixed(3)} MB</td>
                    </tr>
                )}
                </tbody>
            </table>}
        />
    )
}

export default Saves;