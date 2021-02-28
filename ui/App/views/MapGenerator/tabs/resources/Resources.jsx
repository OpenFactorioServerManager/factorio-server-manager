import React from "react";
import ResourceConfigurator from "./components/ResourceConfigurator";

const Resources = () => {
    return <table className="w-full">
        <thead>
            <tr>
                <th colSpan={2}>Name</th>
                <th>Frequency</th>
                <th>Size</th>
                <th>Richness</th>
            </tr>
        </thead>
        <tbody>
            <ResourceConfigurator name="iron"/>
            <ResourceConfigurator name="copper"/>
            <ResourceConfigurator name="coal"/>
            <ResourceConfigurator name="uranium"/>
            <ResourceConfigurator name="crude-oil"/>
        </tbody>
    </table>
}

export default Resources;