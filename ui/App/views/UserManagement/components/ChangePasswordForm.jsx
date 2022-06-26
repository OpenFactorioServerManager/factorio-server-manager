import {useForm} from "react-hook-form";
import React from "react";
import user from "../../../../api/resources/user";
import Button from "../../../components/Button";
import Label from "../../../components/Label";
import Input from "../../../components/Input";
import Error from "../../../components/Error";

const ChangePasswordForm = () => {
    const {register, handleSubmit, formState: {errors}, watch} = useForm();

    const onSubmit = async (data) => {
        const res = await user.changePassword(data);
        if (res) {
            // Update successful
            window.flash("Password changed", "green")
        }
    }

    return (
        <form onSubmit={handleSubmit(onSubmit)}>
            <div className="mb-4">
                <Label htmlFor="old_password" text="Old Password"/>
                <Input register={register('old_password',{required: true})}
                       type="password"
                       placeholder="Old Password"
                />
                <Error error={errors.old_password} message="Old Password is required"/>
            </div>
            <div className="mb-4">
                <Label htmlFor="new_password" text="New Password"/>
                <Input register={register('new_password',{required: true})}
                       type="password"
                       placeholder="New Password"
                />
                <Error error={errors.new_password} message="New Password is required"/>
            </div>
            <div className="mb-4">
                <Label htmlFor="new_password_confirmation" text="New Password Confirmation"/>
                <Input register={register('new_password_confirmation',{required: true})}
                       type="password"
                       placeholder="New Password"
                />
                <Error error={errors.new_password_confirmation} message="New Password Confirmation is required"/>
            </div>
            <Button isSubmit={true} type="success">Change</Button>
        </form>
    )
}

export default ChangePasswordForm
