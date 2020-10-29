import React from "react";
import {faTrashAlt} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import modsResource from "../../../../api/resources/mods";

const ModPack = ({modPack, reloadModPacks}) => {

    const deleteModPack = name => {
        modsResource.packs.delete(name)
            .then(reloadModPacks)
    }

    return (
        <tr>
            <td>{modPack.name}</td>
            <td>
                <FontAwesomeIcon className={"text-red cursor-pointer hover:text-red-light"}
                                 onClick={() => deleteModPack(modPack.name)} icon={faTrashAlt}/>
            </td>
        </tr>
    )
}

export default ModPack;