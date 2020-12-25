import React, {useCallback} from "react";

const TabTitle = ({ title, setSelectedTab, index, isActive }) => {

    const onClick = useCallback(() => {
        setSelectedTab(index)
    }, [setSelectedTab, index])

    return (
            <span className={"accentuated-t accentuated-x cursor-pointer px-3 rounded-t py-1 font-bold relative " + (isActive ? "z-20 text-dirty-white bg-gray-dark" : "z-0 text-black bg-gray-light")} onClick={onClick}>{title}</span>
    )
}

export default TabTitle;