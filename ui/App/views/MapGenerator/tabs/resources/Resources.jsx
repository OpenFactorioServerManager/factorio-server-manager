import React from "react";
import ResourceConfigurator from "./components/ResourceConfigurator";

const Resources = ({settings, setSettings}) => {
    return <table className="flex-1">
        <thead>
            <tr>
                <th colSpan={2}>Name</th>
                <th>Frequency</th>
                <th>Size</th>
                <th>Richness</th>
            </tr>
        </thead>
        <tbody>
            <ResourceConfigurator label="Iron Ore" namePrefix="iron-ore"/>
            <ResourceConfigurator label="Copper Ore" namePrefix="copper-ore"/>
            <ResourceConfigurator label="Stone Ore" namePrefix="stone"/>
            <ResourceConfigurator label="Coal" namePrefix="coal"/>
            <ResourceConfigurator label="Uranium Ore" namePrefix="uranium-ore"/>
            <ResourceConfigurator label="Crude Oil" namePrefix="crude-oil"/>
        </tbody>
    </table>
}

export default Resources;