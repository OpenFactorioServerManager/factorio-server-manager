import Panel from "../../components/Panel";
import React, {useCallback, useEffect, useState} from "react";
import user from "../../../api/resources/user";
import CreateUserForm from "./components/CreateUserForm";
import ChangePasswordForm from "./components/ChangePasswordForm"
import {faTrashAlt} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import { useTranslation } from "react-i18next";

const UserManagement = () => {

    const { t, i18n } = useTranslation();

    const [users, setUsers] = useState([]);

    const updateList = useCallback(async () => {
        const res = await user.list();
        if (res) {
            setUsers(res);
        }
    }, []);

    const deleteUser = useCallback(async (username) => {
        user.delete(username)
            .then(updateList);
    }, []);

    useEffect(() => {
        updateList()
    }, []);

    return (
        <>
            <Panel
                title={t("users.user_list")}
                content={
                    <table className="w-full">
                        <thead>
                        <tr className="text-left py-1">
                            <th>{t("name")}</th>
                            <th>{t("role")}</th>
                            <th>{t("email")}</th>
                            <th>{t("actions")}</th>
                        </tr>
                        </thead>
                        <tbody>
                        {users.map(user =>
                            <tr className="py-1" key={user.username}>
                                <td className="pr-4">{user.username}</td>
                                <td className="pr-4">{user.role}</td>
                                <td className="pr-4">{user.email}</td>
                                <td>
                                    <FontAwesomeIcon className="text-red cursor-pointer hover:text-red-light mr-2" onClick={() => deleteUser(user.username)} icon={faTrashAlt}/>
                                </td>
                            </tr>
                        )}
                        </tbody>
                    </table>
                }
                className="mb-4"
            />
            <Panel
                title={t("users.change_password.title")}
                content={<ChangePasswordForm/>}
                className="mb-4"
            />
            <Panel
                title={t("users.create_user.title")}
                content={<CreateUserForm updateUserList={updateList}/>}
                className="mb-4"
            />
        </>
    )
}

export default UserManagement;