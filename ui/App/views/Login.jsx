import React from 'react';
import {useForm} from "react-hook-form";
import user from "../../api/resources/user";

const Login = () => {
    const {register, handleSubmit, errors} = useForm();

    // todo call api
    const onSubmit = async data => {
        const res = await user.login(data)
        console.log(res)
    };

    return (
        <div className="h-screen overflow-hidden flex items-center justify-center bg-banner">
            <form onSubmit={handleSubmit(onSubmit)} className="rounded-sm bg-gray-dark shadow-xl">
                <div className="px-4 py-2 text-xl text-dirty-white font-bold">
                    Login
                </div>
                <div className="rounded-sm bg-gray-medium shadow-inner mx-4 px-6 pt-4 pb-6 mb-4 flex flex-col">
                    <div className="mb-4">
                        <label className="block text-white text-sm font-bold mb-2" htmlFor="username">
                            Username
                        </label>
                        <input className="shadow appearance-none border w-full py-2 px-3 text-black"
                               ref={register({required: true})}
                               id="username"
                               name="username"
                               type="text" placeholder="Username"/>
                        {errors.password && <span className="block text-red">Username is required</span>}
                    </div>
                    <div className="mb-6">
                        <label className="block text-white text-sm font-bold mb-2" htmlFor="password">
                            Password
                        </label>
                        <input
                            className="shadow appearance-none w-full py-2 px-3 text-black"
                            ref={register({required: true})}
                            name="password"
                            id="password" type="password" placeholder="******************"/>
                        {errors.password && <span className="block text-red">Password is required</span>}
                    </div>
                    <div className="text-center">
                        <button className="bg-green hover:bg-green-light text-black font-bold py-2 px-4 w-full"
                                type="submit">
                            Sign In
                        </button>
                        <a className="bg-gray-light hover:bg-orange text-black py-2 px-4 mt-2 w-full block align-baseline font-bold text-sm text-blue hover:text-blue-darker"
                           href="#">
                            Forgot Password?
                        </a>
                    </div>
                </div>
            </form>
        </div>
    );
};

export default Login;