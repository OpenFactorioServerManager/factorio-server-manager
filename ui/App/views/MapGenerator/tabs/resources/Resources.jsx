import React from "react";
import ResourceConfigurator from "./components/ResourceConfigurator";

const Resources = ({settings, setSettings}) => {
    return <table className="flex-1">
        <thead>
            <tr>
                <th>Name</th>
                <th>Frequency</th>
                <th>Size</th>
                <th>Richness</th>
            </tr>
        </thead>
        <tbody>
            <ResourceConfigurator settings={settings} setSettings={setSettings} label="Iron Ore" name="iron-ore"/>
            <ResourceConfigurator settings={settings} setSettings={setSettings} label="Copper Ore" name="copper-ore"/>
            <ResourceConfigurator settings={settings} setSettings={setSettings} label="Stone Ore" name="stone"/>
            <ResourceConfigurator settings={settings} setSettings={setSettings} label="Coal" name="coal"/>
            <ResourceConfigurator settings={settings} setSettings={setSettings} label="Uranium Ore" name="uranium-ore"/>
            <ResourceConfigurator settings={settings} setSettings={setSettings} label="Crude Oil" name="crude-oil"/>
        </tbody>
    </table>
}

export default Resources;