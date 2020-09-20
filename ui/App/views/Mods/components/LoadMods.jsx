import React, {useEffect, useState} from "react";
import savesResource from "../../../../api/resources/saves";
import Select from "../../../components/Select";
import Label from "../../../components/Label";
import {useForm} from "react-hook-form";
import Button from "../../../components/Button";
import modsResource from "../../../../api/resources/mods";

const LoadMods = ({refreshMods}) => {

    const [saves, setSaves] = useState([]);
    const {register, handleSubmit} = useForm();

    useEffect(() => {
        (async () => {
            setSaves(await savesResource.list());
        })();
    }, []);

    const loadMods = data => {
        savesResource.mods(data.save)
            .then(({mods}) => {
                modsResource.portal.installMultiple(mods).then(refreshMods)
            })
    }

    return (
        <form onSubmit={handleSubmit(loadMods)}>
            <Label text="Save" htmlFor="save"/>
            <Select name="save" inputRef={register} className="mb-4">
                {saves?.map(save => <option value={save.name} key={save.name}>{save.name}</option>)}
            </Select>
            <Button isSubmit={true}>Load</Button>
        </form>
    )
}

export default LoadMods;