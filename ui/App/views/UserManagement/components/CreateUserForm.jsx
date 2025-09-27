import {useForm} from "react-hook-form";
import React from "react";
import user from "../../../../api/resources/user";
import Button from "../../../components/Button";
import Label from "../../../components/Label";
import Input from "../../../components/Input";
import Error from "../../../components/Error";
import { useTranslation } from "react-i18next";

const CreateUserForm = ({updateUserList}) => {

    const { t, i18n } = useTranslation();

    const roleValue = "admin";

    const {
        register,
        handleSubmit,
        formState: {errors},
        watch
    } = useForm({
        values: {
            role: roleValue,
        }
    });
    const password = watch('password');

    const onSubmit = async (data) => {
        const res = await user.add(data);
        if (res) {
            updateUserList()
        }
    }

    return (
        <form onSubmit={handleSubmit(onSubmit)}>
            <div className="mb-4">
                <Label htmlFor="username" text={t("username")}/>
                <Input register={register('username', {required: true})}
                       type="text"
                       placeholder={t("username")}
                />
                <Error error={errors.username} message={t("users.create_user.username_error")}/>
            </div>
            <div className="mb-4">
                <Label htmlFor="role" text={t("role")}/>
                <Input register={register('role', {required: true})}
                       value={roleValue}
                       disabled={true}
                       placeholder={t("Role")}
                />
                <Error error={errors.role} message={t("users.create_user.role_error")}/>
            </div>
            <div className="mb-4">
                <Label htmlFor="email" text={t("email")}/>
                <Input register={register('email', {required: true})}
                       type="email"
                       placeholder={t("email")}
                />
                <Error error={errors.email} message={t("users.create_user.email_error")}/>
            </div>
            <div className="mb-4">
                <Label htmlFor="password" text={t("password")}/>
                <Input register={register('password', {required: true})}
                       type="password"
                       placeholder={t("password")}
                />
                <Error error={errors.password} message={t("users.create_user.password_error")}/>
            </div>
            <div className="mb-4">
                <Label htmlFor="password_confirmation" text={t("users.create_user.password_confirmation")}/>
                <Input register={register('password_confirmation', {
                            required: true,
                            validate: conformation => conformation === password
                        })}

                       type="password"
                       placeholder={t("users.create_user.password_confirmation")}
                />
                <Error error={errors.password_confirmation}
                       message={t("users.create_user.password_confirmation")}/>
            </div>
            <Button isSubmit={true} type="success">{t("save")}</Button>
        </form>
    )
}

export default CreateUserForm;
