import React from "react";

const Input = ({name, inputRef, placeholder = null, type="text", defaultValue=null}) => {
    return (
        <input
            className="shadow appearance-none border w-full py-2 px-3 text-black"
            placeholder={placeholder}
            ref={inputRef}
            name={name}
            id={name}
            type={type}
            defaultValue={defaultValue}
        />
    )
}

export default Input;