import React, {useState} from "react";
import Button from "../../../components/Button";
import Label from "../../../components/Label";
import {useForm} from "react-hook-form";
import modsResource from "../../../../api/resources/mods";
import { useTranslation } from "react-i18next";

const UploadMod = ({refetchInstalledMods}) => {

    const { t, i18n } = useTranslation();

    const defaultFileName = t("mods.select_file")
    const [fileName, setFileName] = useState(defaultFileName);
    const {register, handleSubmit} = useForm();
    const [isUploading, setIsUploading] = useState(false);

    const onSubmit = (data, e) => {
        setIsUploading(true)
        modsResource.upload(data.mod_file[0])
            .then(refetchInstalledMods)
            .finally(() => {
                e.target.reset()
                setFileName(defaultFileName)
                setIsUploading(false);
            })
    }

    return (
        <form onSubmit={handleSubmit(onSubmit)}>
            <Label text={t("save")} htmlFor="mod_file"/>
            <div className="relative bg-white shadow text-black h-full w-full mb-4">
                <input
                    {...register('mod_file')}
                    className="absolute left-0 top-0 opacity-0 cursor-pointer w-full h-full"
                    onChange={e => setFileName(e.currentTarget.files[0].name)}
                    id="mod_file"
                    type="file"
                    accept="application/zip,.zip,.dat,.json"
                />
                <div className="px-2 py-2">{fileName}</div>
            </div>
            <Button isLoading={isUploading} isSubmit={true}>{t("upload")}</Button>
        </form>
    )
}

export default UploadMod;
