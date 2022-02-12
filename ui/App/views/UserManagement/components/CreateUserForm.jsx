import {useForm} from "react-hook-form";
import React from "react";
import user from "../../../../api/resources/user";
import Button from "../../../components/Button";
import Label from "../../../components/Label";
import Input from "../../../components/Input";
import Error from "../../../components/Error";

const CreateUserForm = ({updateUserList}) => {

    const {register, handleSubmit, errors, watch} = useForm();
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
                <Label htmlFor="username" text="Username"/>
                <Input inputRef={register({required: true})}
                       name="username"
                       type="text"
                       placeholder="Username"
                />
                <Error error={errors.username} message="Username is required"/>
            </div>
            <div className="mb-4">
                <Label htmlFor="role" text="Role" />
                <input className="shadow appearance-none border w-full py-2 px-3 text-black"
                       ref={register({required: true})}
                       id="role"
                       name="role"
                       value="admin"
                       disabled={true}
                       type="text"
                       placeholder="Role"
                />
                <Error error={errors.role} message="Role is required"/>
            </div>
            <div className="mb-4">
                <Label htmlFor="email" text="Email"/>
                <Input inputRef={register({required: true})}
                       name="email"
                       type="email"
                       placeholder="Email"
                />
                <Error error={errors.email} message="Email is required"/>
            </div>
            <div className="mb-4">
                <Label htmlFor="password" text="Password"/>
                <Input inputRef={register({required: true})}
                       name="password"
                       type="password"
                       placeholder="Password"
                />
                <Error error={errors.password} message="Password is required"/>
            </div>
            <div className="mb-4">
                <Label htmlFor="password_confirmation" text="Password Confirmation"/>
                <Input inputRef={register({required: true, validate: conformation => conformation === password})}
                       name="password_confirmation"
                       type="password"
                       placeholder="Password Confirmation"
                />
                <Error error={errors.password_confirmation} message="Password Confirmation is required and must match the Password"/>
            </div>
            <Button isSubmit={true} type="success">Save</Button>
        </form>
    )
}

export default CreateUserForm;