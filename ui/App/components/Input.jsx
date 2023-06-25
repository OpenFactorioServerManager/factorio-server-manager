import React from "react";

const Input = ({
                   register,
                   placeholder = null,
                   type = "text",
                   defaultValue = null,
                   hasAutoComplete = true,
                   onKeyDown = () => null,
                   min = null,
                   value = undefined,
                   readOnly = false,
                   disabled = false
               }) => {
    return (
        <input
            className="shadow appearance-none border w-full py-2 px-3 text-black"
            placeholder={placeholder}
            {...register}
            type={type}
            onKeyDown={onKeyDown}
            autoComplete={hasAutoComplete ? "on" : "off"}
            defaultValue={defaultValue}
            min={min}
            value={value}
            readonly={readOnly}
            disabled={disabled}
        />
    )
}

export default Input;
