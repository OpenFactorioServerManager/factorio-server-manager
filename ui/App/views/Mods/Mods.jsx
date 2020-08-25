import Panel from "../../components/Panel";
import React, {useEffect, useState} from "react";
import modsResource from "../../../api/resources/mods";
import Button from "../../components/Button";
import Mod from "./components/Mod";
import server from "../../../api/resources/server";

const Mods = () => {

    const [installedMods, setInstalledMods] = useState([]);
    const [factorioVersion, setFactorioVersion] = useState(null);


    const fetchInstalledMods = () => {
        modsResource.installed()
            .then(res => {
                if (res) {
                    setInstalledMods(res || []);
                }
            });
    };

    useEffect(() => {
        server.factorioVersion()
            .then(res => {
                if (res) {
                    setFactorioVersion(res.base_mod_version)
                }
                fetchInstalledMods();
            })
    }, []);

    return (
        <>
            <Panel
                title="Mods"
                className="mb-6"
                content={
                    <table className="w-full">
                        <thead>
                        <tr className="text-left py-1">
                            <th>Name</th>
                            <th>Enabled</th>
                            <th>Compatibility</th>
                            <th>Mod Version</th>
                            <th>Factorio Version</th>
                            <th/>
                        </tr>
                        </thead>
                        <tbody>
                        {factorioVersion !== null && installedMods.map(mod => <Mod mod={mod} key={mod.name}
                                                                          refreshInstalledMods={fetchInstalledMods}
                                                                          factorioVersion={factorioVersion}/>)}
                        </tbody>
                    </table>

                }

            />
        </>
    )
}

export default Mods;