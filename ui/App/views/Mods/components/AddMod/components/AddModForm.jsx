import React from "react";
import modsResource from "../../../../../../api/resources/mods";
import Button from "../../../../../components/Button";
import Label from "../../../../../components/Label";
import {useForm} from "react-hook-form";
import Input from "../../../../../components/Input";


const AddModForm = ({setIsFactorioAuthenticated}) => {

    const {register, handleSubmit} = useForm();

    const logout = () => {
        modsResource.portal.logout()
            .then(res => {
                setIsFactorioAuthenticated(false);
            });
    }

    return (
        <form >
            <div className="mb-4">
                <Label text="Mod" htmlFor="mod"/>
                <Input inputRef={register({required: true})} name="mod"/>
            </div>
            <Button isSubmit={true}>Install</Button>
            <Button onClick={logout} type="danger">Logout</Button>
        </form>

)
}

export default AddModForm;