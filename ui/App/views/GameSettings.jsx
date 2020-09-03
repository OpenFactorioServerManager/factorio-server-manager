import Panel from "../components/Panel";
import React, {useEffect, useState} from "react";
import settingsResource from "../../api/resources/settings";

const GameSettings = () => {

    const [settingsCategories, setSettingsCategories] = useState();

    const fetchSettings = async () => {
        const res = await settingsResource.game.list();
        setSettingsCategories(res);
    }

    useEffect(() => {
        fetchSettings();
    }, []);

    return (
        <Panel
            className="mb-4"
            title="Game Settings"
            content={
                <>
                    {settingsCategories && Object.keys(settingsCategories).map(key => {
                        const settings = settingsCategories[key];
                        return (
                            <div key={key}>
                                <h1 className="mb-1 text-lg text-dirty-white">{key}</h1>
                                <table key={key} className="w-full mb-2">
                                    <tbody>
                                    {settings && (Object.keys(settings).length > 0 && Object.keys(settings).map(key => {
                                        return (
                                            <tr className="py-1" key={key}>
                                                <td className="w-1/3 pr-4">{key}</td>
                                                <td className="w-2/3 pr-4">{settings[key]}</td>
                                            </tr>
                                        )
                                    })) || <tr>
                                        <td colSpan={2}>--</td>
                                    </tr>}
                                    </tbody>
                                </table>
                            </div>
                        )
                    })}
                </>
            }
        />
    )
}

export default GameSettings;