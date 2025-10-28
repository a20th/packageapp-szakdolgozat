/** @typedef {typeof __propDef.props}  AlertProps */
/** @typedef {typeof __propDef.events}  AlertEvents */
/** @typedef {typeof __propDef.slots}  AlertSlots */
export default class Alert extends SvelteComponent<{
    [x: string]: any;
    heading?: string | undefined;
    class?: string | undefined;
    content?: any;
    closeAriaLabel?: string | undefined;
    closeClassName?: string | undefined;
    color?: string | undefined;
    dismissible?: boolean | undefined;
    fade?: boolean | undefined;
    isOpen?: boolean | undefined;
    toggle?: Function | undefined;
    theme?: string | undefined;
    transition?: object | undefined;
}, {
    [evt: string]: CustomEvent<any>;
}, {
    heading: {};
    default: {};
}> {
}
export type AlertProps = typeof __propDef.props;
export type AlertEvents = typeof __propDef.events;
export type AlertSlots = typeof __propDef.slots;
import { SvelteComponent } from "svelte";
declare const __propDef: {
    props: {
        [x: string]: any;
        heading?: string | undefined;
        class?: string | undefined;
        content?: any;
        closeAriaLabel?: string | undefined;
        closeClassName?: string | undefined;
        color?: string | undefined;
        dismissible?: boolean | undefined;
        fade?: boolean | undefined;
        isOpen?: boolean | undefined;
        toggle?: Function | undefined;
        theme?: string | undefined;
        transition?: object | undefined;
    };
    events: {
        [evt: string]: CustomEvent<any>;
    };
    slots: {
        heading: {};
        default: {};
    };
};
export {};
