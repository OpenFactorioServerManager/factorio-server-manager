import React, {useState} from "react";
import TabTitle from "./TabTitle";

const TabControl = ({children}) => {
    const [selectedTab, setSelectedTab] = useState(0)

    return (
        <div className="mb-6">
            <div className="px-4 pt-3">
                {children.map((item, index) => (
                    <TabTitle
                        key={index}
                        title={item.props.title}
                        index={index}
                        isActive={index === selectedTab}
                        setSelectedTab={setSelectedTab}
                    />
                ))}
            </div>
            <div className="z-10 relative accentuated bg-gray-dark p-4">
                <div className="text-white rounded-sm bg-gray-medium shadow-inner px-6 pt-4 pb-6">
                    {children[selectedTab]}
                </div>
            </div>
        </div>
    )
}

export default TabControl