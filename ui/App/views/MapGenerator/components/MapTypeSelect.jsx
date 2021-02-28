import React, {useState} from "react";

const MapTypeSelect = ({inputRef, onChange}) => {

    const [value, setValue] = useState("default");

    const options = [
        {
            name: "Default",
            value: "default"
        },
        {
            name: "Rich resources",
            value: "rich-resources"
        },
        {
            name: "Marathon",
            value: "marathon"
        },
        {
            name: "Death world",
            value: "death-world"
        },
        {
            name: "Death world marathon",
            value: "death-world-marathon"
        },
        {
            name: "Rail world",
            value: "rail-world"
        },
        {
            name: "Ribbon world",
            value: "ribbon-world"
        },
        {
            name: "Island",
            value: "island"
        },
    ];

    const change = optionElement => {
        setValue(optionElement.target.value)
        // onChange(optionElement)
    }

    return <div className="relative ">
        <select
            className="shadow appearance-none bg-gray-light w-full h-8 px-1 text-black"
            name={name}
            ref={inputRef}
            id={name}
            value={value}
            onChange={change}
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
}

export default MapTypeSelect;