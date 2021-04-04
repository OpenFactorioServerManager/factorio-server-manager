import React, {useState, useEffect} from "react";
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
import saves from "../../../api/resources/saves";
import MapPreviewImage from "./components/MapPreviewImage";


const MapGenerator = () => {

    const {register, handleSubmit} = useForm();
    const [seed, setSeed] = useState(0);
    const [settings, setSettings] = useState({});
    const [previewImage, setPreviewImage] = useState(null);
    const [isLoadingPreview, setIsLoadingPreview] = useState(false);

    const loadPreview = () => {
        setIsLoadingPreview(true)

        saves.preview(settings)
            .then(imageData => setPreviewImage(imageData))
            .finally(() => setIsLoadingPreview(false))
    }

    const randomSeed = () => {
        const randomValue = Math.floor(Math.random() * 1000000000)
        updateSeed(randomValue)
    }

    const updateSeed = value => {
        setSeed(value)
        setSettings(Object.assign(settings, {seed: value}))
    }

    useEffect(() => {
        console.log('test')

        Promise.all([
            saves.defaultMapGenSettings()
                .then(mapGenSettings => setSettings(Object.assign(settings,mapGenSettings))),
            saves.defaultMapSettings()
                .then(mapSettings => setSettings(Object.assign(settings,mapSettings))),

        ])
            .finally(() => {
                randomSeed()
                loadPreview()
            })

    }, []);

    return <form onSubmit={handleSubmit(data => null)}>
        <TabControl
            actions={
                <Button size="sm" isSubmit={true} onClick={loadPreview} isLoading={isLoadingPreview} type="success">Generate Map</Button>
            }
            title={
                <div className="flex justify-between my-1">
                    <MapTypeSelect/>
                    <SeedInput updateSeed={updateSeed} seed={seed} generateRandomSeed={randomSeed}/>
                </div>
            }
        >
            <Tab title="Resources">
                <div className="flex">
                    <Resources/>
                    <MapPreviewImage imageData={previewImage}/>
                </div>
            </Tab>
            <Tab title="Terrain">
                <div className="flex">
                    <Terrain/>
                    <MapPreviewImage imageData={previewImage}/>
                </div>
            </Tab>
            <Tab title="Enemy">
                <div className="flex">
                    <Enemy/>
                    <MapPreviewImage imageData={previewImage}/>
                </div>
            </Tab>
            <Tab title="Advanced">
                <div className="flex">
                    <Advanced/>
                    <MapPreviewImage imageData={previewImage}/>
                </div>
            </Tab>
        </TabControl>
    </form>
}

export default MapGenerator;