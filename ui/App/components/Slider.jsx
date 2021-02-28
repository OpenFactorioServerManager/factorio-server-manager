import React from "react";

const Slider = () => {
    return <input type="range" className="slider" defaultValue={6} min={1} max={12} step={1}/>
}

export default Slider;