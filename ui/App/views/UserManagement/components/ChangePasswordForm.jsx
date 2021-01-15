import {useForm} from "react-hook-form";
import React from "react";
import user from "../../../../api/resources/user";
import Button from "../../../components/Button";

const ChangePasswordForm = () => {
    const {register, handleSubmit, errors, watch} = useForm();
    const password = watch('new_password');

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
                <label className="block text-white text-sm font-bold mb-2" htmlFor="old_password">
                    Old Password:
                </label>
                <input className="shadow appearance-none border w-full py-2 px-3 text-black"
                       ref={register({required: true})}
                       id="old_password"
                       name="old_password"
                       type="password"
                       placeholder="**********"
                />
                {errors.old_password && <span className="block text-red">Old Password is required</span>}
            </div>
            <div className="mb-4">
                <label className="block text-white text-sm font-bold mb-2" htmlFor="new_password">
                    New Password:
                </label>
                <input className="shadow appearance-none border w-full py-2 px-3 text-black"
                       ref={register({required: true})}
                       id="new_password"
                       name="new_password"
                       type="password"
                       placeholder="**********"
                />
                {errors.new_password && <span className="block text-red">New Password is required</span>}
            </div>
            <div className="mb-4">
                <label className="block text-white text-sm font-bold mb-2" htmlFor="new_password_confirmation">
                    New Password Confirmation:
                </label>
                <input className="shadow appearance-none border w-full py-2 px-3 text-black"
                       ref={register({required: true, validate: confirmation => confirmation === password})}
                       id="new_password_confirmation"
                       name="new_password_confirmation"
                       type="password"
                       placeholder="**********"
                />
                {errors.new_password_confirmation && <span className="block text-red">New Password Confirmation is required</span>}
            </div>
            <Button isSubmit={true} type="success">Change</Button>
        </form>
    )
}

export default ChangePasswordForm