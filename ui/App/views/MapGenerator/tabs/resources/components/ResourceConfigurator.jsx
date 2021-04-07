import React from "react";
import Input from "../../../../../components/Input";

const ResourceConfigurator = ({label, name, settings, setSettings }) => {

    const resource = settings?.autoplace_controls?.[name]

    const updateResource = (attribute, value) => {
        let tmpSettings = JSON.parse(JSON.stringify(settings));
        tmpSettings.autoplace_controls[name][attribute] = parseFloat(value)
        setSettings(tmpSettings);
    }

    return <tr >
        <td className="p-2">
             {label}
        </td>
        <td className="p-2">
            <Input className="text-center" value={resource?.frequency ? resource.frequency : ""} onChange={event => {
                updateResource('frequency', event.target.value)
            }}/>
        </td>
        <td className="p-2">
            <Input className="text-center" value={resource?.size ? resource.size : ""} onChange={event => {
                updateResource('size', event.target.value)
            }}/>
        </td>
        <td className="p-2">
            <Input className="text-center" value={resource?.richness ? resource.richness : ""} onChange={event => {
                updateResource('richness', event.target.value)
            }}/>
        </td>
    </tr>
}

export default ResourceConfigurator;