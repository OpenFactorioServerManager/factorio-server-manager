import React, {useState} from 'react';
import Modal from "./Modal";
import Button from "./Button";
import { useTranslation } from "react-i18next";

function ChangeLangDialog({isOpen, close, onSuccess}) {

    const { t, i18n } = useTranslation();

    const changeLang = (langKey) => {
        i18n.changeLanguage(langKey)
        close()
    }

    return (
        <Modal
            title={t("lang")}
            content={
                <>
                    <Button className="w-full" onClick={() => changeLang("en")}>English</Button>
                    <Button className="w-full" onClick={() => changeLang("ru")}>Русский</Button>
                    {/* Other languages */}
                </>
            }
            actions={
            <>
                <Button size="sm" type="danger" className="mr-2" onClick={close}>{t("cancel")}</Button>
            </>}
            isOpen={isOpen}
            onSuccess={onSuccess}
        />
    );
}

export default ChangeLangDialog;