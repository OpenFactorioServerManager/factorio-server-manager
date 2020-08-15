import React from "react";

const Select = ({name, inputRef, children}) => {
    return (
        <div className="relative">
        <select
            className="shadow appearance-none border w-full py-2 px-3 text-black"
            name={name}
            ref={inputRef}
            id={name}
        >
            {children.map(child => child)}
        </select>
        <div
            className="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-black">
            <svg className="fill-current h-4 w-4" xmlns="http://www.w3.org/2000/svg"
                 viewBox="0 0 20 20">
                <path
                    d="M9.293 12.95l.707.707L15.657 8l-1.414-1.414L10 10.828 5.757 6.586 4.343 8z"/>
            </svg>
        </div>
        </div>
    )
}

export default Select;