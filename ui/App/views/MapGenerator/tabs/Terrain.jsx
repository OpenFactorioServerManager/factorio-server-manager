import React from "react";
import Select from "../../../components/Select";
import Checkbox from "../../../components/Checkbox";
import Input from "../../../components/Input";
import Label from "../../../components/Label";
import copy from "../../../copy";

const Terrain = ({settings, setSettings}) => {
    return <div className="flex-1">
        <div className="flex justify-between border border-black rounded p-2 mb-1">
            <div className="self-center">
                <Label text="Maptype"/>
            </div>
            <Select isInline options={[
                {value: 'normal', name: 'Normal'},
                {value: 'island', name: 'Island'}
            ]}/>
        </div>
        <div className="border border-black rounded p-2 mb-1">
            <table className="w-full table-fixed">
                <thead>
                <tr>
                    <th className="w-3/5"/>
                    <th className="w-1/5">Scale</th>
                    <th className="w-1/5">Coverage</th>
                </tr>
                </thead>
                <tbody>
                <tr>
                    <td><Checkbox checked={true} text="Water"/></td>
                    <td><Input/></td>
                    <td><Input/></td>
                </tr>
                <tr>
                    <td><Checkbox checked={true} text="Trees"/></td>
                    <td><Input/></td>
                    <td><Input/></td>
                </tr>
                </tbody>
            </table>
        </div>
        <div className="border border-black rounded p-2 mb-1">
            <table className="w-full table-fixed">
                <thead>
                <tr>
                    <th className="w-3/5"/>
                    <th className="w-1/5">Frequency</th>
                    <th className="w-1/5">Continuity</th>
                </tr>
                </thead>
                <tbody>
                <tr>
                    <td>
                        <Checkbox text="Cliffs"
                                  checked
                        />
                    </td>
                    <td>
                        <Input
                            value={40 / (settings?.cliff_settings?.cliff_elevation_interval  || 40)}
                            onChange={event => {
                                let tmp = copy(settings);
                                tmp.cliff_settings.cliff_elevation_interval = 40 / parseInt(event.target.value);
                                setSettings(tmp);
                            }}
                        />
                    </td>
                    <td>
                        <Input
                            value={settings?.cliff_settings?.richness || 1}
                            onChange={event => {
                                let tmp = copy(settings);
                                tmp.cliff_settings.richness = parseInt(event.target.value);
                                setSettings(tmp);
                            }}
                        />
                    </td>
                </tr>
                </tbody>
            </table>
        </div>
        <div className="border border-black rounded p-2">
            <table className="w-full table-fixed">
                <thead>
                <tr>
                    <th className="w-3/5"/>
                    <th className="w-1/5">Scale</th>
                    <th className="w-1/5">Bias</th>
                </tr>
                </thead>
                <tbody>
                <tr>
                    <td>Moisture</td>
                    <td><Input/></td>
                    <td><Input/></td>
                </tr>
                <tr>
                    <td>Terrain Typ</td>
                    <td><Input/></td>
                    <td><Input/></td>
                </tr>
                </tbody>
            </table>
        </div>
    </div>
}

export default Terrain;