import React from "react";
import {faSpinner} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";


const MapPreviewImage = ({imageData, show, isLoading, seed}) => {
    return <>
        {show && <div className="flex-1 accentuated rounded-sm relative">
            {isLoading &&
                <>
                    <div className="absolute z-20 opacity-50 bg-black w-full h-full"/>
                    <div className="absolute z-30 flex justify-center items-center w-full h-full">
                        <div className="text-center opacity-100">
                            <FontAwesomeIcon className="inline-block" size="4x" icon={faSpinner} spin={true}/>
                            <p>Loading Map Preview</p>
                        </div>
                    </div>
                </>
            }
            <div className="absolute z-10 bottom-0 right-0">
                <p className="text-xs px-1">Seed: {seed}</p>
            </div>
            <img src={imageData}/>
        </div>}
    </>
}

export default MapPreviewImage;