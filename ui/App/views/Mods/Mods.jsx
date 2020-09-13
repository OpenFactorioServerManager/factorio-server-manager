import Panel from "../../components/Panel";
import React, {useEffect, useState} from "react";
import modsResource from "../../../api/resources/mods";
import Button from "../../components/Button";
import Mod from "./components/Mod";
import server from "../../../api/resources/server";
import TabControl from "../../components/Tabs/TabControl";
import Tab from "../../components/Tabs/Tab";
import AddMod from "./components/AddMod/AddMod";
import UploadMod from "./components/UploadMod";
import LoadMods from "./components/LoadMods";

const Mods = () => {

    const [installedMods, setInstalledMods] = useState([]);
    const [factorioVersion, setFactorioVersion] = useState(null);


    const fetchInstalledMods = () => {
        modsResource.installed()
            .then(setInstalledMods);
    };

    const deleteAllMods = () => {
        modsResource.deleteAll()
            .then(fetchInstalledMods)
    }

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
        <div>
            <TabControl>
                <Tab title="Install Mod">
                    <AddMod refetchInstalledMods={fetchInstalledMods}/>
                </Tab>
                <Tab title="Upload Mod">
                    <UploadMod/>
                </Tab>
                <Tab title="Load Mod from Save">
                    <LoadMods/>
                </Tab>
            </TabControl>

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
                actions={
                    <Button size="sm" type="danger" onClick={deleteAllMods}>Delete all Mods</Button>
                }
            />

            <Panel
                title="Mod packs"
                className="mb-6"
                content={
                   "Test"
                }
            />
        </div>
    )
}

export default Mods;