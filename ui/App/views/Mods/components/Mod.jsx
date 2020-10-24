import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {
    faArrowCircleUp,
    faCheck,
    faSpinner,
    faTimes,
    faToggleOff,
    faToggleOn,
    faTrashAlt
} from "@fortawesome/free-solid-svg-icons";
import modsResource from "../../../../api/resources/mods";
import React, {useEffect, useState} from "react";
import {coerce, gt, satisfies} from "semver";

const Mod = ({refreshInstalledMods, mod, factorioVersion}) => {

    const [newVersion, setNewVersion] = useState(null)
    const [icon, setIcon] = useState(faArrowCircleUp)

    const toggleMod = async modName => {
        modsResource
            .toggle(modName)
            .then(refreshInstalledMods)
    }

    const deleteMod = async modName => {
        modsResource
            .delete(modName)
            .then(refreshInstalledMods)
    }

    const updateMod = async (modName, download_url, file_name) => {
        setIcon(faSpinner);
        modsResource.update(modName, download_url, file_name)
            .then(async () => {
                await refreshInstalledMods();
            })
            .finally(() => setIcon(faArrowCircleUp))
    }

    useEffect(() => {
        (async () => {
            const data = await modsResource.portal.info(mod.name)

            //get newest COMPATIBLE release
            let newestRelease;
            data.releases.forEach(release => {
                if (
                    gt(coerce(release.version).version, coerce(mod.version).version) && (
                        satisfies(factorioVersion, coerce(release.info_json.factorio_version).version) ||
                        (satisfies(factorioVersion, "1.0.0") && satisfies(coerce(release.info_json.factorio_version).version, "0.18.x"))
                    )
                ) {
                    if (!newestRelease) {
                        newestRelease = release;
                    } else if (gt(coerce(release.version).version, coerce(newestRelease.version).version)) {
                        newestRelease = release;
                    }
                }
            });

            if (newestRelease && newestRelease.version !== mod.version) {
                setNewVersion({
                    downloadUrl: newestRelease.download_url,
                    file_name: newestRelease.file_name,
                    version: newestRelease.version
                });
            } else {
                setNewVersion(null);
            }

        })();
    }, [mod]);

    return (
        <tr className="py-1">
            <td className="pr-4">{mod.title}</td>
            <td className="pr-4">
                {mod.enabled
                    ? <FontAwesomeIcon className="cursor-pointer hover:text-green-light text-green" icon={faToggleOn}
                                       onClick={() => toggleMod(mod.name)}/>
                    : <FontAwesomeIcon className="cursor-pointer hover:text-red-light text-red" icon={faToggleOff}
                                       onClick={() => toggleMod(mod.name)}/>
                }
            </td>
            <td className="pr-4">
                {mod.compatibility
                    ? <FontAwesomeIcon className="text-green" icon={faCheck}/>
                    : <FontAwesomeIcon className="text-red" icon={faTimes}/>
                }
            </td>
            <td className="pr-4">{mod.version} {newVersion && <FontAwesomeIcon spin={icon === faSpinner}
                                                                               onClick={() => updateMod(mod.name, newVersion.downloadUrl, mod.file_name)}
                                                                               className="hover:text-orange cursor-pointer ml-1"
                                                                               icon={icon}/>}</td>
            <td className="pr-4">{mod.factorio_version}</td>
            <td className="pr-4">
                <FontAwesomeIcon className={"text-red cursor-pointer hover:text-red-light"}
                                 onClick={() => deleteMod(mod.name)} icon={faTrashAlt}/>
            </td>
        </tr>
    )
}

export default Mod;

