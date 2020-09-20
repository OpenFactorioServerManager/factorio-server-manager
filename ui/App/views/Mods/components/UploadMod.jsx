import React, {useState} from "react";
import Button from "../../../components/Button";
import Label from "../../../components/Label";

const UploadMod = () => {

    const [fileName, setFileName] = useState('Select File ...');

    return (
        <form>
            <Label text="Save" htmlFor="savefile"/>
            <div className="relative bg-white shadow text-black h-full w-full mb-4">
                <input
                    className="absolute left-0 top-0 opacity-0 cursor-pointer w-full h-full"
                    onChange={e => setFileName(e.currentTarget.files[0].name)}
                    name="savefile"
                    id="savefile" type="file"/>
                <div className="px-2 py-2">{fileName}</div>
            </div>
            <Button isSubmit={true}>Upload</Button>
        </form>
    )
}

export default UploadMod;