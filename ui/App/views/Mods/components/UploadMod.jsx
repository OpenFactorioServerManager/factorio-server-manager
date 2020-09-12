import React, {useState} from "react";

const UploadMod = () => {

    const [fileName, setFileName] = useState('Select File ...');

    return (
        <div className="relative bg-white shadow text-black h-full w-full">
            <input
                className="absolute left-0 top-0 opacity-0 cursor-pointer w-full h-full"
                onChange={e => setFileName(e.currentTarget.files[0].name)}
                name="savefile"
                id="savefile" type="file"/>
            <div className="px-2 py-3">{fileName}</div>
        </div>
    )
}

export default UploadMod;