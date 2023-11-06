import Panel from "../components/Panel";
import React from "react";

const Help = () => {
    return (
        <Panel
            title="Help"
            content={
                <>
                    <h1 className="text-xl text-dirty-white">Factorio Server Manager</h1>
                    <p className="mb-2">The Factorio Server Manager (FSM) is an open source project and is not affiliated to the game Factorio or Wube Software.</p>

                    <h2 className="text-dirty-white">Bugs and Help</h2>
                    <p className="mb-4">Please use the <a className="text-blue hover:text-blue-light" target="_blank" href="https://github.com/OpenFactorioServerManager/factorio-server-manager/issues">GitHub repository</a> to report bugs or seek for help.</p>

                    <h1 className="mb-1 text-xl text-dirty-white">Helpful Resources</h1>
                    <p className="mb-2"><a className="text-blue hover:text-blue-light" target="_blank" href="https://wiki.factorio.com/Multiplayer">Official Factorio Wiki about Multiplayer</a></p>
                </>
            }
        />
    )
}

export default Help;