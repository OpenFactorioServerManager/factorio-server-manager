import ModsContent from "../ModsContent.jsx";

function instanceOfModsContentFunction(isRequired) {
    return function(props, propName, componentName) {
        if(props[propName]) {
            if(!props[propName] instanceof ModsContent) {
                return new Error(propName + ' in ' + componentName + ' is not an instance of ModContent');
            }
        } else {
            if(isRequired) {
                return new Error(propName + ' in ' + componentName + ' is missing');
            }
        }

        return null;
    }
}

const instanceOfModsContent = instanceOfModsContentFunction(false);
instanceOfModsContent.isRequired = instanceOfModsContentFunction(true);

export {instanceOfModsContent};