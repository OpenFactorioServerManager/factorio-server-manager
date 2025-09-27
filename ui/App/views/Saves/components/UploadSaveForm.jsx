import Button from "../../../components/Button";
import React, {useState} from "react";
import {useForm} from "react-hook-form";
import saves from "../../../../api/resources/saves";
import Error from "../../../components/Error";
import { useTranslation } from "react-i18next";


const UploadSaveForm = ({onSuccess}) => {
    const { t, i18n } = useTranslation();
    const {register, handleSubmit, formState: {errors}} = useForm();
    const [fileName, setFileName] = useState(t("saves.upload_form.select_file"));

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
                    {t("saves.upload_form.file_name")}
                </label>
                <div className="relative bg-white shadow text-black w-full">
                    <input
                        className="absolute left-0 top-0 opacity-0 cursor-pointer w-full h-full"
                        {...register("savefile", {required: true})}
                        onChange={e => setFileName(e.currentTarget.files[0].name)}
                        type="file"/>
                    <div className="px-2 py-3">{fileName}</div>
                </div>
                <Error error={errors.savefile} message={t("saves.upload_form.save_file_error_message")}/>
            </div>
            <Button type="success" isSubmit={true}>{t("saves.upload_form.upload")}</Button>
        </form>
    )
}

export default UploadSaveForm;
