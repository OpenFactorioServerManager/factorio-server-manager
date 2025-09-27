import Panel from "../components/Panel";
import React from "react";
import { useTranslation, Trans } from "react-i18next";

const Help = () => {

    const { t, i18n } = useTranslation();
    
    return (
        <Panel
            title={t("help.title")}
            content={
                <>
                    <h1 className="text-xl text-dirty-white">{t("help.fsm")}</h1>
                    <p className="mb-2">{t("help.fsm_content")}</p>

                    <h2 className="text-dirty-white">{t("help.bugs_help")}</h2>
                    <p className="mb-4">
                        <Trans
                            i18nKey="help.bugs_help_content"
                            components={[
                            <a
                                className="text-blue hover:text-blue-light"
                                target="_blank"
                                href="https://github.com/OpenFactorioServerManager/factorio-server-manager/issues"
                            />
                            ]}
                        />
                    </p>

                    <h1 className="mb-1 text-xl text-dirty-white">{t("help.helpful_resources")}</h1>
                    <p className="mb-2"><a className="text-blue hover:text-blue-light" target="_blank" href="https://wiki.factorio.com/Multiplayer">{t("help.factorio_link_text")}</a></p>
                </>
            }
        />
    )
}

export default Help;