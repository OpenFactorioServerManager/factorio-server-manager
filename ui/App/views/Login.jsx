import React, {useEffect} from 'react';
import {useForm} from "react-hook-form";
import user from "../../api/resources/user";
import Button from "../components/Button";
import {useHistory, useLocation} from "react-router";
import Panel from "../components/Panel";
import Input from "../components/Input";
import Label from "../components/Label";
import {Flash} from "../components/Flash";
import Error from "../components/Error";

const Login = ({handleLogin}) => {
    const {register, handleSubmit, errors} = useForm();
    const history = useHistory();
    const location = useLocation();

    const onSubmit = async data => {
        try {
            const loginAttempt = await user.login(data)
            if (loginAttempt?.username) {
                await handleLogin(loginAttempt);
                history.push('/');
            }
        } catch (e) {
            console.log(e);
            window.flash("Login failed. Username or Password wrong.", "red");
            throw e;
        }
    };

    // on mount check if user is authenticated
    useEffect(() => {
        (async () => {
            const status = await user.status();
            if (status?.username) {
                await handleLogin(status);
                history.push(location?.state?.from || '/');
            }
        })();
    }, [])

    return (
        <div className="h-screen overflow-hidden flex items-center justify-center bg-black">
            <Panel
                title="Login"
                content={
                    <form onSubmit={handleSubmit(onSubmit)}>
                        <div className="mb-4">
                            <Label text="Username" htmlFor="username"/>
                            <Input inputRef={register({required: true})} name="username" placeholder="Username"/>
                            <Error error={errors.username} message="Username is required"/>
                        </div>
                        <div className="mb-6">
                            <Label text="Password" htmlFor="password"/>
                            <Input
                                inputRef={register({required: true})}
                                name="password"
                                type="password"
                                placeholder="******************"
                            />
                            <Error error={errors.password} message="Password is required"/>
                        </div>
                        <div className="text-center">
                            <Button type="success" className="w-full" isSubmit={true}>Sign In</Button>
                        </div>
                    </form>
                }
            />
            <Flash/>
        </div>
    );
};

export default Login;