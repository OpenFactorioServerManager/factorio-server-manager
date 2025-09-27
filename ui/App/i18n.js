import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import Backend from 'i18next-http-backend'
import LanguageDetector from "i18next-browser-languagedetector";
import enTranslations from "./locales/en.json";
import ruTranslations from "./locales/ru.json";

const resources =
{
    en:
    {
        translation: enTranslations,
    },
    ru:
    {
        translation: ruTranslations,
    }
};

i18n
    .use(Backend)
    .use(LanguageDetector)
    .use(initReactI18next)
    .init({
        resources,
        fallbackLng: "en", // Default Language
        // Detecting and caching of language cookies
        detection:
        {
            order: ["localStorage", "cookie", "navigator"],
            cache: ["localStorage", "cookie"]
        },
        interpolation:
        {
            escapeValue: false // react already safes from xss
        }
    });

export default i18n;