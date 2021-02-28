import React from "react";
import TabControl from "../../components/Tabs/TabControl";
import Tab from "../../components/Tabs/Tab";
import Button from "../../components/Button";
import Resources from "./tabs/resources/Resources";
import Terrain from "./tabs/Terrain";
import Enemy from "./tabs/Enemy";
import Advanced from "./tabs/Advanced";
import {useForm} from "react-hook-form";
import SeedInput from "./components/SeedInput";
import MapTypeSelect from "./components/MapTypeSelect";

const MapGenerator = () => {

    const {register, handleSubmit} = useForm()

    return <form onSubmit={handleSubmit(data => null)}>
        <TabControl
            actions={
                <Button size="sm" isSubmit={true} type="success">Generate Map</Button>
            }
            title={
                <div className="flex justify-between my-1">
                    <MapTypeSelect/>
                    <SeedInput inputRef={register}/>
                </div>
            }
        >
            <Tab title="Resources">
                <Resources/>
            </Tab>
            <Tab title="Terrain">
                <Terrain/>
            </Tab>
            <Tab title="Enemy">
                <Enemy/>
            </Tab>
            <Tab title="Advanced">
                <Advanced/>
            </Tab>
        </TabControl>
    </form>
}

export default MapGenerator;