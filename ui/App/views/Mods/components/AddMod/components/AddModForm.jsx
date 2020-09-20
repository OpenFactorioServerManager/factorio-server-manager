import React, {useCallback, useEffect, useState} from "react";
import modsResource from "../../../../../../api/resources/mods";
import Button from "../../../../../components/Button";
import Label from "../../../../../components/Label";
import {useForm} from "react-hook-form";
import Input from "../../../../../components/Input";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faExternalLinkAlt} from "@fortawesome/free-solid-svg-icons/faExternalLinkAlt";
import {faSpinner} from "@fortawesome/free-solid-svg-icons";

const LinkModPortal = () => {
    return <a href="https://mods.factorio.com" target="_blank" className="px-2 text-blue hover:text-blue-light">Mod
        Portal <FontAwesomeIcon icon={faExternalLinkAlt}/></a>
}

const AddModForm = ({setIsFactorioAuthenticated, fuse}) => {

    const {register, watch, setValue} = useForm();
    const [suggestedMods, setSuggestedMods] = useState([]);
    const [selectedMod, setSelectedMod] = useState(null)

    const [autocomplete, setAutocomplete] = useState(NaN);
    const mod = watch('mod');

    const logout = () => {
        modsResource.portal.logout()
            .then(() => setIsFactorioAuthenticated(false));
    }

    const updateSuggestedMods = () => {
        if (typeof fuse != "undefined") {
            clearTimeout(autocomplete)
            setAutocomplete(setTimeout(() => setSuggestedMods(fuse.search(mod || '')), 200));
        }
    };

    const selectMod = mod => {
        clearTimeout(autocomplete);
        setValue('mod', mod.item.title); // triggers effect for mod
        setSuggestedMods([]);
        setSelectedMod(mod);
    }

    // get triggered if mod changed
    useEffect(() => {
        if (selectedMod === null) {
            updateSuggestedMods();
        } else if (selectedMod.item.title !== mod) {
            setSelectedMod(null);
            updateSuggestedMods();
        }
    }, [mod]);

    const addMod = data => {
        // todo install selected mod
        console.log(data);
        // todo update list of installed mods
    }

    return (
        <form onSubmit={addMod}>
            <div className="mb-4 relative">
                <Label text="Mod" htmlFor="mod"/>
                { typeof fuse !== "undefined"
                    ? <Input inputRef={register({required: true})} hasAutoComplete={false} name="mod"/>
                    : <div className="border border-gray-medium w-full py-2 px-3 text-white">
                        <FontAwesomeIcon icon={faSpinner} spin={true}/> Loading List of Mods from <LinkModPortal/>
                    </div>
                }
                {suggestedMods.length > 0 &&
                    <ul className="bg-white text-black h-64 overflow-y-scroll absolute bottom-0 left-0 w-full -mb-64">
                        {suggestedMods.map((mod, index) => <li className="px-2 py-1 cursor-pointer hover:bg-blue-light" onClick={() => selectMod(mod)} key={index}>{mod.item.title}</li>)}
                    </ul>
                }
            </div>
            <Button isDisabled={selectedMod === null} isSubmit={true} className="mr-2">Install</Button>
            <Button onClick={logout} type="danger" className="mr-2">Logout</Button>
            <LinkModPortal/>
        </form>
    )
}

export default AddModForm;