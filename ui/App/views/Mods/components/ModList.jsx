import Mod from "./Mod";
import React from "react";
import { useTranslation } from "react-i18next";


const ModList = ({mods, factorioVersion, updateMod, toggleMod, deleteMod, addUpdatableMod = null, disabled = false}) => {

    const { t, i18n } = useTranslation();

    return (
        <table className="w-full">
            <thead>
            <tr className="text-left py-1">
                <th>{t("name")}</th>
                <th>{t("mods.mod_list.enabled")}</th>
                <th>{t("mods.mod_list.compatibility")}</th>
                <th>{t("mods.mod_list.mod_version")}</th>
                <th>{t("mods.mod_list.factorio_version")}</th>
                <th/>
            </tr>
            </thead>
            <tbody>
            {
                factorioVersion !== null && mods.map(
                    (mod, i) =>
                        <Mod mod={mod} key={i}
                             updateMod={updateMod}
                             toggleMod={toggleMod}
                             deleteMod={deleteMod}
                             addUpdatableMod={addUpdatableMod}
                             factorioVersion={factorioVersion}
                             disabled={disabled}
                        />
                )
            }
            </tbody>
        </table>
    )
}

export default ModList;