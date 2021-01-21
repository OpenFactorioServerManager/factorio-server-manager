import React from "react"

const Error = ({error, message}) => {
    if (error) {
        return (
            <span className="block text-red">
                {message}
            </span>
        )
    }
    return null
}

export default Error