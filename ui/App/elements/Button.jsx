import React from "react";

const Button = ({ children, type, onClick, isSubmit }) => {

    let color = '';

    switch (type) {
        case 'success':
            color = 'bg-green hover:bg-green-light';
            break;
        case 'danger':
            color = 'bg-red hover:bg-red-light';
            break;
        default:
            color = 'bg-gray-light hover:bg-orange'
    }

    return (
        <button onClick={onClick} className={color + '  text-black font-bold py-2 px-4 w-full'}
                type={isSubmit ? 'submit' : 'button'}>
            {children}
        </button>
    );
}

export default Button;