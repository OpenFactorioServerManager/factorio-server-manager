import React, {useState} from 'react';
import Modal from "./Modal";
import Button from "./Button";
import { useTranslation } from "react-i18next";

function ConfirmDialog({title, content, isOpen, close, onSuccess}) {

    const { t, i18n } = useTranslation();

    const [isLoading, setIsLoading] = useState(false);

    const confirm = () => {
        setIsLoading(true);
        onSuccess()
            .finally(() => {
                close();
                setIsLoading(false);
            })
    }

    return (
        <Modal
            title={title}
            content={content}
            actions={
                <>
                    <Button size="sm" type="danger" className="mr-2" onClick={close}>{t("cancel")}</Button>
                    <Button size="sm" isLoading={isLoading} type="success" onClick={confirm}>{t("confirm")}</Button>
                </>
            }
            isOpen={isOpen}
        />
    );
}

export default ConfirmDialog;