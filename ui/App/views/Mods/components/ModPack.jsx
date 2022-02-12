import React, {useState} from "react";
import {faSpinner, faTrashAlt, faUpload} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import modsResource from "../../../../api/resources/mods";
import ModList from "./ModList";
import ConfirmDialog from "../../../components/ConfirmDialog";

const ModPack = ({modPack, reloadModPacks, factorioVersion, reloadMods, disabled = false}) => {

    const [isLoading, setIsLoading] = useState(false);
    const [isLoadModPackDialogOpen, setIsLoadModPackDialogOpen] = useState(false);


    const deleteModPack = modName => {
        return modsResource.packs
            .delete(modName)
            .then(reloadModPacks)
    }

    const toggleMod = modName => {
        return modsResource
            .packs
            .mods
            .toggle(modPack.name, modName)
            .then(reloadModPacks)
    }

    const updateMod = version => {
        return modsResource
            .packs
            .mods
            .update(modPack.name, version)
            .then(reloadModPacks)
    }

    const deleteMod = modName => {
        return modsResource
            .packs
            .mods
            .delete(modPack.name, modName)
            .then(reloadModPacks)
    }

    const loadModPack = name => {
        setIsLoading(true)
        return modsResource.packs
            .load(name)
            .then(reloadMods)
            .finally(() => setIsLoading(false))
    }

    return (
        <div className="mb-4">
            <div className="flex items-center justify-between">
                <h2 className="text-lg text-dirty-white mb-1 inline">{modPack.name}</h2>
                <div className="flex space-x-2">
                    {
                        !disabled &&
                        <>
                            <FontAwesomeIcon className="text-blue cursor-pointer hover:text-blue-light inline"
                                             onClick={() => setIsLoadModPackDialogOpen(true)}
                                             spin={isLoading}
                                             icon={isLoading ? faSpinner : faUpload}
                            />
                            <ConfirmDialog
                                title="Load ModPack"
                                content={`Loading the ModPack ${modPack.name} will remove all installed Mods.`}
                                isOpen={isLoadModPackDialogOpen}
                                close={() => setIsLoadModPackDialogOpen(false)}
                                onSuccess={() => loadModPack(modPack.name)}
                            />
                        </>
                    }

                    <FontAwesomeIcon className="text-red cursor-pointer hover:text-red-light inline"
                                     onClick={() => deleteModPack(modPack.name)} icon={faTrashAlt}/>
                </div>
            </div>
            <ModList mods={modPack.mods.mods}
                     factorioVersion={factorioVersion}
                     toggleMod={toggleMod}
                     updateMod={updateMod}
                     deleteMod={deleteMod}
            />
        </div>
    )
}

export default ModPack;