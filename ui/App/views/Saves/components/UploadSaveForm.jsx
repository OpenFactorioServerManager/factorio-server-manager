import Button from "../../../components/Button";
import React, {useState} from "react";
import {useForm} from "react-hook-form";
import saves from "../../../../api/resources/saves";
import Error from "../../../components/Error";
import Label from "../../../components/Label";


const UploadSaveForm = ({onSuccess}) => {
    const {register, handleSubmit, errors} = useForm();
    const [fileName, setFileName] = useState('Select File ...');

    const onSubmit = (data, e) => {
        saves.upload(data.savefile[0]).then(_ => {
            e.target.reset();
            onSuccess();
        })
    };

    return (
        <form onSubmit={handleSubmit(onSubmit)}>
            <div className="mb-6">
                <Label text="Save File" htmlFor="savefile"/>
                <div className="relative rounded-sm bg-gray-light shadow text-black w-full">
                    <input
                        className="absolute  left-0 top-0 opacity-0 cursor-pointer w-full h-full"
                        ref={register({required: true})}
                        onChange={e => setFileName(e.currentTarget.files[0].name)}
                        name="savefile"
                        id="savefile" type="file"/>
                    <div className="px-3 py-2">{fileName}</div>
                </div>
                <Error error={errors.savefile} message="Savefile is required"/>
            </div>
            <Button type="success" isSubmit={true}>Upload</Button>
        </form>
    )
}

export default UploadSaveForm;