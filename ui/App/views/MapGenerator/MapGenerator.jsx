import React, {useEffect, useState, useCallback} from "react";
import TabControl from "../../components/Tabs/TabControl";
import Tab from "../../components/Tabs/Tab";
import Button from "../../components/Button";
import Resources from "./tabs/resources/Resources";
import Terrain from "./tabs/Terrain";
import Enemy from "./tabs/Enemy";
import Advanced from "./tabs/Advanced";
import SeedInput from "./components/SeedInput";
import MapTypeSelect from "./components/MapTypeSelect";
import saves from "../../../api/resources/saves";
import MapPreviewImage from "./components/MapPreviewImage";

let timeoutPreviewHandle = null;

const MapGenerator = () => {

    const [isPreviewDisplayed, setIsPreviewDisplayed] = useState(false)
    const [seed, setSeed] = useState(0);
    const [settings, setSettings] = useState({});
    const [previewImage, setPreviewImage] = useState(null);
    const [previewImageSeed, setPreviewImageSeed] = useState(null);
    const [isLoadingPreview, setIsLoadingPreview] = useState(false);

    const updateSettings = newSettings => {
        setSettings(newSettings)
        loadPreview();
    }


    const loadPreview = (force = false) => {
        if (isLoadingPreview || (!isPreviewDisplayed && !force)) {
            return;
        }

        const previewHandler = () => {
            setIsLoadingPreview(true)
            setPreviewImageSeed(settings.seed)

            saves.preview(settings)
                .then(imageData => setPreviewImage(imageData))
                .finally(() => setIsLoadingPreview(false))
        }

        if (timeoutPreviewHandle) {
            clearTimeout(timeoutPreviewHandle);
            timeoutPreviewHandle = null
        }
        timeoutPreviewHandle = setTimeout(previewHandler, 600);
    }

    const randomSeed = () => {
        const randomValue = Math.floor(Math.random() * 1000000000)
        updateSeed(randomValue)
    }

    const updateSeed = value => {
        setSeed(value);
        setSettings(Object.assign(settings, {seed: value}));
        loadPreview();
    }

    useEffect(() => {
        Promise.all([
            saves.defaultMapGenSettings()
                .then(mapGenSettings => setSettings(Object.assign(settings, mapGenSettings))),
            saves.defaultMapSettings()
                .then(mapSettings => setSettings(Object.assign(settings, mapSettings))),

        ]).finally(() => {
            randomSeed()
            loadPreview()
        })
    }, []);

    return <TabControl
        actions={
            <div className="flex justify-between">
                <Button size="sm"  type="success">Generate Map</Button>
                {isPreviewDisplayed
                    ? <Button size="sm" onClick={() => setIsPreviewDisplayed(false)}>Hide Preview</Button>
                    : <Button size="sm" onClick={() => {
                        setIsPreviewDisplayed(true)
                        loadPreview(true)
                    }
                    }>Show Preview</Button>
                }
            </div>
        }
        title={
            <div className="flex justify-between my-1">
                <MapTypeSelect/>
                <SeedInput updateSeed={updateSeed} seed={seed} generateRandomSeed={randomSeed}/>
            </div>
        }
    >
        <Tab title="Resources">
            <div className="flex space-x-8">
                <Resources settings={settings} setSettings={updateSettings}/>
                <MapPreviewImage imageData={previewImage} isLoading={isLoadingPreview} show={isPreviewDisplayed} seed={previewImageSeed}/>
            </div>
        </Tab>
        <Tab title="Terrain">
            <div className="flex space-x-8">
                <Terrain settings={settings} setSettings={updateSettings}/>
                <MapPreviewImage imageData={previewImage} isLoading={isLoadingPreview} show={isPreviewDisplayed} seed={previewImageSeed}/>
            </div>
        </Tab>
        <Tab title="Enemy">
            <div className="flex space-x-8">
                <Enemy settings={settings} setSettings={updateSettings}/>
                <MapPreviewImage imageData={previewImage} isLoading={isLoadingPreview} show={isPreviewDisplayed} seed={previewImageSeed}/>
            </div>
        </Tab>
        <Tab title="Advanced">
            <div className="flex space-x-8">
                <Advanced settings={settings} setSettings={updateSettings}/>
                <MapPreviewImage imageData={previewImage} isLoading={isLoadingPreview} show={isPreviewDisplayed} seed={previewImageSeed}/>
            </div>
        </Tab>
    </TabControl>
}

export default MapGenerator;