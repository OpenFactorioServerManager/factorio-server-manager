import React from "react";
import Label from "../../../components/Label";
import Input from "../../../components/Input";

const Advanced = () => {
    return <div className="flex justify-between w-1/2">
        <div>
            Map Size
        </div>
        <div>
            <Label text="Height"/>
            <Input/>
            <Label text="Width"/>
            <Input/>
        </div>
    </div>
}

export default Advanced;