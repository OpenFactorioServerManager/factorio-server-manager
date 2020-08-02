import {useForm} from "react-hook-form";
import React from "react";
import user from "../../../../api/resources/user";
import Button from "../../../components/Button";

const CreateUserForm = ({updateUserList}) => {

    const {register, handleSubmit, errors, watch} = useForm();
    const password = watch('password');

    const onSubmit = async (data) => {
        const res = user.add(data);
        if (res.success) {
            updateUserList()
        }
    }

    return (
        <form onSubmit={handleSubmit(onSubmit)}>
            <div className="mb-4">
                <label className="block text-white text-sm font-bold mb-2" htmlFor="username">
                    Username
                </label>
                <input className="shadow appearance-none border w-full py-2 px-3 text-black"
                       ref={register({required: true})}
                       id="username"
                       name="username"
                       type="text" placeholder="Username"/>
                {errors.username && <span className="block text-red">Username is required</span>}
            </div>
            <div className="mb-4">
                <label className="block text-white text-sm font-bold mb-2" htmlFor="username">
                    Role
                </label>
                <input className="shadow appearance-none border w-full py-2 px-3 text-black"
                       ref={register({required: true})}
                       id="role"
                       name="role"
                       value="admin"
                       disabled={true}
                       type="text" placeholder="Role"/>
                {errors.role && <span className="block text-red">Role is required</span>}
            </div>
            <div className="mb-4">
                <label className="block text-white text-sm font-bold mb-2" htmlFor="username">
                    Email
                </label>
                <input className="shadow appearance-none border w-full py-2 px-3 text-black"
                       ref={register({required: true})}
                       id="email"
                       name="email"
                       type="text" placeholder="Email"/>
                {errors.email && <span className="block text-red">Email is required</span>}
            </div>
            <div className="mb-4">
                <label className="block text-white text-sm font-bold mb-2" htmlFor="username">
                    Password
                </label>
                <input className="shadow appearance-none border w-full py-2 px-3 text-black"
                       ref={register({required: true})}
                       id="password"
                       name="password"
                       type="password" placeholder="Password"/>
                {errors.password && <span className="block text-red">Password is required</span>}
            </div>
            <div className="mb-4">
                <label className="block text-white text-sm font-bold mb-2" htmlFor="username">
                    Password Confirmation
                </label>
                <input className="shadow appearance-none border w-full py-2 px-3 text-black"
                       ref={register({required: true, validate: confirmation => confirmation === password})}
                       id="password_confirmation"
                       name="password_confirmation"
                       type="password" placeholder="Password Confirmation"/>
                {errors.password_confirmation && <span className="block text-red">Password Confirmation is required and must match the Password</span>}
            </div>
            <Button isSubmit={true} type="success">Save</Button>
        </form>
    )
}

export default CreateUserForm;