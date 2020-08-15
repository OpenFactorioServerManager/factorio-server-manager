import React from "react";

const Label = ({text, htmlFor}) => {
    return (
        <label
            className="block text-white text-sm font-bold mb-1" htmlFor={htmlFor}>
            {text}
        </label>
    )
}

export default Label;