import React, {useEffect, useState} from "react";
import savesResource from "../../../../api/resources/saves";
import Select from "../../../components/Select";
import Label from "../../../components/Label";
import {useForm} from "react-hook-form";
import Button from "../../../components/Button";

const LoadMods = () => {

    const [saves, setSaves] = useState([]);
    const {register, handleSubmit} = useForm();

    useEffect(() => {
        (async () => {
            const res = await savesResource.list()
            if (res.success) {
                setSaves(res.data);
            }
        })();
    }, [])

    return (
        <form >
            <Label text="Save" htmlFor="save"/>
            <Select name="save" inputRef={register}>
                {saves?.map(save => <option value={save.name} key={save.name}>{save.name}</option>)}
            </Select>
            <Button isSubmit={true}>Save</Button>
        </form>
    )
}

export default LoadMods;