import Mod from "./Mod";
import React from "react";


const ModList = ({mods, factorioVersion, updateMod, toggleMod, deleteMod, addUpdatableMod = null, disabled = false}) => {

    return (
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