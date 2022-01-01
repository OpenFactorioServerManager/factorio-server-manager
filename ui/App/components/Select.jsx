import React, {useState} from "react";

const Select = ({register, options, className = "", defaultValue = ""}) => {

    const [value, setValue] = useState(defaultValue);

    return (
        <div className={`${className} relative`}>
        <select
            className="shadow appearance-none border w-full py-2 px-3 text-black"
            name={name}
            {...register}
            value={value}
            onChange={optionElement => setValue(optionElement.target.value)}
        >
            {options.map(option => <option value={option.value} key={option.value}>{option.name}</option>)}
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
