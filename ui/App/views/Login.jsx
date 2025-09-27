import React, {useEffect} from 'react';
import {useForm} from "react-hook-form";
import user from "../../api/resources/user";
import Button from "../components/Button";
import {useLocation, useNavigate} from "react-router";
import Panel from "../components/Panel";
import Input from "../components/Input";
import Label from "../components/Label";
import {Flash} from "../components/Flash";
import Error from "../components/Error";
import { useTranslation } from "react-i18next";

const Login = ({handleLogin}) => {

    const { t, i18n } = useTranslation();

    const {register, handleSubmit, formState: { errors }} = useForm();
    const navigate = useNavigate();
    const location = useLocation();

    const onSubmit = async data => {
        try {
            const loginAttempt = await user.login(data)
            if (loginAttempt?.username) {
                await handleLogin(loginAttempt);
                navigate('/');
            }
        } catch (e) {
            console.log(e);
            window.flash(t("login.login_failed_message"), "red");
            throw e;
        }
    };

    // on mount check if user is authenticated
    useEffect(() => {
        (async () => {
            const status = await user.status();
            if (status?.username) {
                await handleLogin(status);
                navigate(location?.state?.from || '/');
            }
        })();
    }, [])

    return (
        <div className="h-screen overflow-hidden flex items-center justify-center bg-black">
            <Panel
                title={t("login.title")}
                content={
                    <form onSubmit={handleSubmit(onSubmit)}>
                        <div className="mb-4">
                            <Label text={t("username")} htmlFor="username"/>
                            <Input register={register('username', {required: true})} placeholder="Username"/>
                            <Error error={errors.username} message={t("login.username_error_message")}/>
                        </div>
                        <div className="mb-6">
                            <Label text={t("password")} htmlFor="password"/>
                            <Input
                                register={register('password',{required: true})}
                                type="password"
                                placeholder="******************"
                            />
                            <Error error={errors.password} message={t("login.password_error_message")}/>
                        </div>
                        <div className="text-center">
                            <Button type="success" className="w-full" isSubmit={true}>{t("login.sign_in")}</Button>
                        </div>
                    </form>
                }
            />
            <Flash/>
        </div>
    );
};

export default Login;
