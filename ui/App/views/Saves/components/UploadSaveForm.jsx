import Button from "../../../components/Button";
import React, {useState} from "react";
import {useForm} from "react-hook-form";
import saves from "../../../../api/resources/saves";
import Error from "../../../components/Error";


const UploadSaveForm = ({onSuccess}) => {
    const {register, handleSubmit, formState: {errors}} = useForm();
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
                <label className="block text-white text-sm font-bold mb-2" htmlFor="password">
                    Save File
                </label>
                <div className="relative bg-white shadow text-black w-full">
                    <input
                        className="absolute left-0 top-0 opacity-0 cursor-pointer w-full h-full"
                        register={register('savefile', {required: true})}
                        onChange={e => setFileName(e.currentTarget.files[0].name)}
                        type="file"/>
                    <div className="px-2 py-3">{fileName}</div>
                </div>
                <Error error={errors.savefile} message="Savefile is required"/>
            </div>
            <Button type="success" isSubmit={true}>Upload</Button>
        </form>
    )
}

export default UploadSaveForm;
