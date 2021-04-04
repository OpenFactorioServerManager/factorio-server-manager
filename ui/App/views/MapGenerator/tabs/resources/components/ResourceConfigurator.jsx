import React from "react";
import Slider from "../../../../../components/Slider";

const ResourceConfigurator = ({ label, namePrefix }) => {
    return <tr>
        <td>
            {/*<input*/}
            {/*    type="checkbox"*/}
            {/*    checked={true}*/}
            {/*/>*/}
        </td>
        <td className="px-2">
             {label}
        </td>
        <td className="px-2">
            <Slider/>
        </td>
        <td className="px-2">
            <Slider/>
        </td>
        <td className="px-2">
            <Slider/>
        </td>
    </tr>
}

export default ResourceConfigurator;