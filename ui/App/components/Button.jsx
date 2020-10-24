import React from "react";
import {faSpinner} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";

const Button = ({ children, type, onClick, isSubmit, className, size, isLoading, isDisabled = false }) => {

    let color = '';
    let padding = '';

    switch (type) {
        case 'success':
            color = `bg-green ${isDisabled ? null : "hover:glow-green hover:bg-green-light" }`;
            break;
        case 'danger':
            color = `bg-red ${isDisabled ? null : "hover:glow-red hover:bg-red-light"}`;
            break;
        default:
            color = `bg-gray-light ${isDisabled ? null : "hover:glow-orange hover:bg-orange"}`
    }

    switch (size) {
        case 'sm':
            padding = 'py-1 px-2';
            break;
        default:
            padding = 'py-2 px-4'
    }

    return (
        <button onClick={onClick} disabled={isDisabled} className={`${className ? className: null} ${isDisabled ? "bg-opacity-50 cursor-not-allowed" : null} ${padding} ${color} inline-block accentuated text-black font-bold`}
        type={isSubmit ? 'submit' : 'button'}>
            {children} { isLoading && <FontAwesomeIcon icon={faSpinner} spin={true}/>}
        </button>
    );
}

export default Button;