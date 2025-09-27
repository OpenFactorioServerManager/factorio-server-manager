import {useForm} from "react-hook-form";
import React from "react";
import user from "../../../../api/resources/user";
import Button from "../../../components/Button";
import Label from "../../../components/Label";
import Input from "../../../components/Input";
import Error from "../../../components/Error";
import { useTranslation } from "react-i18next";

const ChangePasswordForm = () => {

    const { t, i18n } = useTranslation();

    const {register, handleSubmit, reset, formState: {errors}, watch} = useForm();

    const new_password = watch("new_password");

    const onSubmit = async (data) => {
        const res = await user.changePassword(data);
        if (res) {
            // Update successful
            window.flash(t("users.change_password.update_successful"), "green")
            reset();
        }
    }

    return (
        <form onSubmit={handleSubmit(onSubmit)}>
            <div className="mb-4">
                <Label htmlFor="old_password" text={t("users.change_password.old_password")}/>
                <Input register={register("old_password",{required: true})}
                       type="password"
                       placeholder={t("users.change_password.old_password")}
                />
                <Error error={errors.old_password} message={t("users.change_password.old_password_error")}/>
            </div>
            <div className="mb-4">
                <Label htmlFor="new_password" text={t("users.change_password.new_password")}/>
                <Input register={register("new_password",{required: true})}
                       type="password"
                       placeholder={t("users.change_password.new_password")}
                />
                <Error error={errors.new_password} message={t("users.change_password.new_password_error")}/>
            </div>
            <div className="mb-4">
                <Label htmlFor="new_password_confirmation" text={t("users.change_password.new_password_confirmation")}/>
                <Input register={register("new_password_confirmation",{required: true, validate: value => value === new_password})}
                       type="password"
                       placeholder={t("users.change_password.new_password")}
                />
                <Error error={errors.new_password_confirmation} message={t("users.change_password.new_password_confirmation_error")}/>
            </div>
            <Button isSubmit={true} type="success">{t("change")}</Button>
        </form>
    )
}

export default ChangePasswordForm
