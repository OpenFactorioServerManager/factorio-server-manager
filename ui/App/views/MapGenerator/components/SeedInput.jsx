import React, {useState, useEffect} from "react";
import Input from "../../../components/Input";
import Button from "../../../components/Button";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faRandom} from "@fortawesome/free-solid-svg-icons";



const SeedInput = ({inputRef}) => {

    const [seed, setSeed] = useState(0);

    const randomSeed = () => {
        setSeed(Math.floor(Math.random() * 1000000000))
    }
    useEffect(randomSeed, []);

    return <div className="relative">
        <div className="w-32 inline-block">
            <input
                className="shadow border-t-2 border-gray-light bg-gray-light h-8 rounded-l-sm appearance-none w-full px-1 text-black"
                ref={inputRef}
                value={seed}
                onChange={setSeed}/>
        </div>

        <Button size="none" onClick={randomSeed} className="text-center ml-1 rounded-r-sm w-8 h-8">
            <FontAwesomeIcon icon={faRandom}/>
        </Button>
    </div>
}

export default SeedInput;