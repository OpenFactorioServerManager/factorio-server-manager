import React from "react";

const Checkbox = ({name, text, inputRef, checked, className, textSize = 'sm', onChange = null}) => {
    return (
        <label className={`block text-gray-500 font-bold ${className}`}>
            <input
                className="mr-2 leading-tight"
                ref={inputRef}
                name={name}
                id={name}
                onChange={onChange}
                type="checkbox" defaultChecked={checked}/>
            <span className={`text-${textSize}`}>{text}</span>
        </label>
    )
}

export default Checkbox;