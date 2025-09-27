import {useForm} from "react-hook-form";
import Button from "../../../components/Button";
import React, {useState} from "react";
import saves from "../../../../api/resources/saves";
import Label from "../../../components/Label";
import Input from "../../../components/Input";
import Error from "../../../components/Error";
import { useTranslation } from "react-i18next";

const CreateSaveForm = ({onSuccess}) => {
    const {register, handleSubmit, formState: {errors}} = useForm();
    const [isLoading, setIsLoading] = useState(false);

    const { t, i18n } = useTranslation();
    

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
                <Label text={t("saves.save_form.file_name")} htmlFor="savefile"/>
                <Input register={register("savefile", {required: true})}/>
                <Error error={errors.savefile} message={t("saves.save_form.save_file_error_message")}/>
            </div>
            <Button type="success" isLoading={isLoading} isSubmit={true}>{t("saves.save_form.create_save")}</Button>
        </form>
    )
}

export default CreateSaveForm;
