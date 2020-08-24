import {useForm} from "react-hook-form";
import Button from "../../../components/Button";
import React from "react";
import saves from "../../../../api/resources/saves";
import Label from "../../../components/Label";
import Input from "../../../components/Input";

const CreateSaveForm = ({onSuccess}) => {
    const {register, handleSubmit, errors} = useForm();


    const onSubmit = async (data, e) => {
        const res = await saves.create(data.savefile);
        if (res) {
            e.target.reset();
            onSuccess();
        }
    };

    return (
        <form onSubmit={handleSubmit(onSubmit)}>
        <div className="mb-6">
            <Label text="Savefile Name" htmlFor="savefile"/>
            <Input inputRef={register({required: true})} name="savefile"/>
            {errors.savefile && <span className="block text-red">Savefile Name is required</span>}
        </div>
        <Button type="success" isSubmit={true}>Create Save</Button>
    </form>
    )
}

export default CreateSaveForm;