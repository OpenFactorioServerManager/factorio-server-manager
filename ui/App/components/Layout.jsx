import React, {useEffect, useState} from "react";
import {NavLink, Outlet} from "react-router-dom";
import Button from "./Button";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faBars} from "@fortawesome/free-solid-svg-icons";
import {Flash} from "./Flash";
import ChangeLangDialog from "./ChangeLangDialog";
import { useTranslation } from "react-i18next";

const Layout = ({handleLogout, serverStatus}) => {

    const { t, i18n } = useTranslation();

    const [isNavCollapsed, setIsNavCollapsed] = useState(true);
    const [isChangingLang, setIsChangingLang] = useState(false);

    const Status = ({info}) => {

        let text = t("controls.unknown");
        let color = 'gray-light';

        if (info && info.running) {
            text = t("controls.running");
            color = 'green';
        } else if (info && !info.running) {
            text = t("controls.stopped");
            color = 'red';
        }

        return (
            <div className={`bg-${color} accentuated rounded px-2 py-1 text-black`}>{text}</div>
        )
    }

    const Link = ({children, to, last}) => {
        return (
            <NavLink
                onClick={() => setIsNavCollapsed(true)}
                end
                to={to}
                className={({isActive}) => {
                    return [
                        isActive ? "bg-orange" : "",
                        `hover:glow-orange accentuated bg-gray-light hover:bg-orange text-black font-bold py-2 px-4 w-full block${last ? '' : ' mb-1'}`,
                    ].join(" ")
                }}
            >{children}</NavLink>)
    }

    return (
        <>
            {/*Sidebar*/}
            <div className="w-full md:w-88 md:fixed md:top-0 md:left-0 bg-gray-dark md:h-screen overflow-y-auto">
                <div className="py-4 px-2 accentuated">
                    <div className="mx-4 justify-between flex text-center">
                        <span className="text-dirty-white text-xl">{t("main_title")}</span>
                        <button
                            className="md:hidden cursor-pointer text-white hover:text-dirty-white"
                            onClick={() => setIsNavCollapsed(!isNavCollapsed)}
                        >
                            <FontAwesomeIcon icon={faBars}/>
                        </button>
                    </div>
                </div>
                <div className={isNavCollapsed ? "hidden md:block" : "block"}>
                    <div className="py-4 px-2 accentuated">
                        <h1 className="text-dirty-white text-lg mb-2 mx-4">{t("server_status")}</h1>
                        <div className="mx-4 mb-4 text-center">
                            <Status info={serverStatus}/>
                        </div>
                    </div>
                    <div className="py-4 px-2 accentuated">
                        <h1 className="text-dirty-white text-lg mb-2 mx-4">{t("server_management")}</h1>
                        <div className="text-white text-center rounded-sm bg-black shadow-inner mx-4 p-1">
                            <Link to="/">{t("controls.title")}</Link>
                            <Link to="/saves">{t("saves.title")}</Link>
                            <Link to="/mods">{t("mods.title")}</Link>
                            <Link to="/server-settings">{t("server_settings.title")}</Link>
                            <Link to="/game-settings">{t("game_settings.title")}</Link>
                            <Link to="/console">{t("console.title")}</Link>
                            <Link to="/logs" last={true}>{t("logs.title")}</Link>
                        </div>
                    </div>
                    <div className="py-4 px-2 accentuated">
                        <h1 className="text-dirty-white text-lg mb-2 mx-4">{t("FSM_administration")}</h1>
                        <div className="text-white text-center rounded-sm bg-black shadow-inner mx-4 p-1">
                            <Link to="/user-management">{t("users.title")}</Link>
                            <Button className="w-full mb-1" onClick={() => setIsChangingLang(true)}>{t("lang")}</Button>
                            <Link to="/help" last={true}>{t("help.title")}</Link>
                            <ChangeLangDialog
                                isOpen={isChangingLang}
                                close={() => setIsChangingLang(false)}
                                onSuccess={() => console.log("new lang apply")}
                            />
                        </div>
                    </div>
                    <div className="py-4 px-2 accentuated">
                        <div className="text-white text-center rounded-sm bg-black shadow-inner mx-4 p-1">
                            <Button type="danger" className="w-full" onClick={handleLogout}>{t("logout")}</Button>
                        </div>
                    </div>
                    <div className="accentuated-t accentuated-x md:block hidden"/>
                </div>
            </div>

            {/*Main*/}
            <div className="md:ml-88 min-h-screen">
                <div className="container md:mx-auto pt-16 md:px-6">
                    <Outlet />
                    <Flash/>
                </div>
            </div>
        </>
    );
}

export default Layout;