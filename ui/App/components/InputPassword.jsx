import Input from "./Input";
import React, {useState} from "react";
import {faEye, faEyeSlash} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";

const InputPassword = ({register, defaultValue}) => {

    const [type, setType] = useState("password");

    let icon;
    if (type === "password") {
        icon = faEye;
    } else {
        icon = faEyeSlash
    }

    return (
        <div className="flex">
            <Input type={type} defaultValue={defaultValue} {...register} placeholder="*************"/>
            <div
                className="accentuated cursor-pointer bg-gray-light flex items-center px-2 text-black"
                onClick={() => setType(type === "password" ? "text" : "password")}
            >
                <FontAwesomeIcon fixedWidth={true} icon={icon} />
            </div>
        </div>
    )
}
export default InputPassword;
