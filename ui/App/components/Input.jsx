import React from "react";

const Input = ({
                   name,
                   inputRef,
                   placeholder = null,
                   type = "text",
                   defaultValue = null,
                   hasAutoComplete = true,
                   onKeyDown = () => null,
                   onChange = () => null,
                   min = null,
                   value = undefined,
                   disabled = false
               }) => {
    return (
        <input
            className="shadow bg-gray-light rounded-sm appearance-none w-full py-2 px-3 text-black"
            placeholder={placeholder}
            ref={inputRef}
            name={name}
            id={name}
            type={type}
            onKeyDown={onKeyDown}
            onChange={onChange}
            autoComplete={hasAutoComplete ? "on" : "off"}
            defaultValue={defaultValue}
            min={min}
            value={value}
            disabled={disabled}
        />
    )
}

export default Input;