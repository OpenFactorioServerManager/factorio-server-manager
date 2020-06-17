import React from "react";

const Button = ({ children, type, onClick, isSubmit, className }) => {

    let color = '';

    switch (type) {
        case 'success':
            color = 'bg-green hover:glow-green hover:bg-green-light';
            break;
        case 'danger':
            color = 'bg-red hover:glow-red hover:bg-red-light';
            break;
        default:
            color = 'bg-gray-light hover:glow-orange hover:bg-orange'
    }

    return (
        <button onClick={onClick} className={`${className} ${color} accentuated text-black font-bold py-2 px-4`}
                type={isSubmit ? 'submit' : 'button'}>
            {children}
        </button>
    );
}

export default Button;