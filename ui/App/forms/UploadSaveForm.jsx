import Button from "../elements/Button";
import React from "react";
import {useForm} from "react-hook-form";
import saves from "../../api/resources/saves";


const UploadSaveForm = ({onSuccess}) => {

    const {register, handleSubmit, errors} = useForm();
    const onSubmit = async data => {
        const res = await saves.upload(data.savefile);
        if (res.success) {
            onSuccess();
        }
    };

    return (
        <form onSubmit={handleSubmit(onSubmit)}>
            <div className="mb-6">
                <label className="block text-white text-sm font-bold mb-2" htmlFor="password">
                    Savefile Name
                </label>
                <input
                    className="shadow appearance-none w-full py-2 px-3 text-black"
                    ref={register({required: true})}
                    name="savefile"
                    id="savefile" type="file"/>
                {errors.savefile && <span className="block text-red">Savefile Name is required</span>}
            </div>
            <Button type="success">Upload Save</Button>
        </form>
    )
}

export default UploadSaveForm;