import Layout from "../components/Layout";
import React, {useEffect, useState} from "react";
import save from "../../api/resources/saves";

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
        <Layout>
            <div className="rounded-sm bg-gray-dark shadow-xl pb-4">
                <div className="px-4 py-2 text-xl text-dirty-white font-bold">
                    Saves
                </div>
                <div className="text-white rounded-sm bg-gray-medium shadow-inner mx-4 px-6 pt-4 pb-6">
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
                            {saves.map(save => (
                                <tr className="py-1" key={save.name}>
                                    <td className="pr-4">{save.name}</td>
                                    <td className="pr-4">{save.last_mod}</td>
                                    <td>{save.size}</td>
                                    <td/>
                                </tr>
                            ))}
                        </tbody>
                    </table>

                </div>
            </div>
        </Layout>
    )
}

export default Saves;