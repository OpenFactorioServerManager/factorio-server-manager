import {useForm} from "react-hook-form";
import Button from "../../../components/Button";
import React, {useState} from "react";
import saves from "../../../../api/resources/saves";
import Label from "../../../components/Label";
import Input from "../../../components/Input";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faSpinner} from "@fortawesome/free-solid-svg-icons";
import Error from "../../../components/Error";

const CreateSaveForm = ({onSuccess}) => {
    const {register, handleSubmit, errors} = useForm();
    const [isLoading, setIsLoading] = useState(false);

    const onSubmit = async (data, e) => {
        setIsLoading(true)
        saves.create(data.savefile)
            .then(() => {
                e.target.reset();
                onSuccess();
            })
            .finally(() => setIsLoading(false))
    };

    return (
        <form onSubmit={handleSubmit(onSubmit)}>
            <div className="mb-6">
                <Label text="Savefile Name" htmlFor="savefile"/>
                <Input inputRef={register({required: true})} name="savefile"/>
                <Error error={errors.savefile} message="Savefile Name is required"/>
            </div>
            <Button type="success" isLoading={isLoading} isSubmit={true}>Create Save</Button>
        </form>
    )
}

export default CreateSaveForm;